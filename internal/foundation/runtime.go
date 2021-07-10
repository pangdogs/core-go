package foundation

import (
	"github.com/pangdogs/core/internal"
	"github.com/pangdogs/core/internal/list"
	"time"
)

type RuntimeWhole interface {
	internal.Runtime
	internal.GC
	AddEntity(entity internal.Entity)
	RemoveEntity(entID uint64)
	PushSafeCall(callBundle *SafeCallBundle)
}

func NewRuntime(ctx internal.Context, app internal.App, optFuncs ...NewRuntimeOptionFunc) internal.Runtime {
	rt := &Runtime{}

	opts := &RuntimeOptions{}
	for _, optFun := range append([]NewRuntimeOptionFunc{NewRuntimeOption.Default()}, optFuncs...) {
		optFun(opts)
	}

	rt.InitRuntime(ctx, app, opts)

	return rt.inheritor
}

type Runtime struct {
	Runnable
	internal.Context
	RuntimeOptions
	id              uint64
	app             internal.App
	safeCallList    chan *SafeCallBundle
	entityMap       map[uint64]*list.Element
	entityList      list.List
	entityStartList []*list.Element
	entityGCList    []*list.Element
	gcList          []internal.GC
	frame           FrameWhole
}

func (rt *Runtime) InitRuntime(ctx internal.Context, app internal.App, opts *RuntimeOptions) {
	if ctx == nil {
		panic("nil ctx")
	}

	if app == nil {
		panic("nil app")
	}

	if opts == nil {
		panic("nil opts")
	}

	rt.RuntimeOptions = *opts

	if rt.inheritor != nil {
		rt.inheritor.(RuntimeInheritorWhole).initRuntimeInheritor(rt)
	} else {
		rt.inheritor = rt
	}

	rt.InitRunnable()
	rt.Context = ctx
	rt.id = app.MakeUID()
	rt.app = app
	rt.safeCallList = make(chan *SafeCallBundle)
	close(rt.safeCallList)
	rt.entityList.Init()
	rt.entityMap = map[uint64]*list.Element{}

	CallOuter(rt.autoRecover, rt.GetReportError(), rt.initFunc)

	if opts.autoRun {
		rt.Run()
	}
}

func (rt *Runtime) Run() chan struct{} {
	if !rt.MarkRunning() {
		panic("runtime already running")
	}

	rt.safeCallList = make(chan *SafeCallBundle, rt.safeCallCacheSize)

	go func() {
		if parentCtx, ok := rt.GetParentContext().(internal.Context); ok {
			parentCtx.GetWaitGroup().Add(1)
		}

		invokeFun := func(fun func(entity EntityWhole)) {
			if fun == nil {
				return
			}
			rt.RangeEntities(func(entity internal.Entity) bool {
				CallOuter(rt.autoRecover, rt.GetReportError(), func() {
					fun(entity.(EntityWhole))
				})
				return true
			})
		}

		startChan := make(chan struct{})

		notifyStart := func() {
			if len(rt.entityStartList) > 0 {
				go func() {
					startChan <- struct{}{}
				}()
			}
		}

		invokeStartFun := func() {
			count := len(rt.entityStartList)
			if count <= 0 {
				return
			}
			for _, e := range rt.entityStartList {
				if e.Escape() || e.GetMark(0) {
					continue
				}
				CallOuter(rt.autoRecover, rt.GetReportError(), e.Value.(EntityWhole).CallStart)
			}
			rt.entityStartList = rt.entityStartList[count:]
		}

		invokeSafeCallFun := func(callBundle *SafeCallBundle) (ret internal.SafeRet) {
			exception := CallOuter(rt.autoRecover, rt.GetReportError(), func() {
				if callBundle.Stack != nil {
					ret = callBundle.SafeFun(callBundle.Stack)
				} else {
					ret = callBundle.UnsafeFun()
				}
			})
			if exception != nil {
				ret.Err = exception
			}
			return
		}

		defer func() {
			close(rt.safeCallList)

			func() {
				for {
					select {
					case <-startChan:
						invokeStartFun()
						notifyStart()
						rt.GC()

					case callBundle, ok := <-rt.safeCallList:
						if !ok {
							return
						}

						callBundle.Ret <- invokeSafeCallFun(callBundle)
						notifyStart()
						rt.GC()

					default:
						return
					}
				}
			}()

			if parentCtx, ok := rt.GetParentContext().(internal.Context); ok {
				parentCtx.GetWaitGroup().Done()
			}

			rt.GetWaitGroup().Wait()
			rt.MarkShutdown()
			rt.shutChan <- struct{}{}

			CallOuter(rt.autoRecover, rt.GetReportError(), rt.stopFunc)
		}()

		rt.frame = nil

		if rt.frameCreatorFunc == nil {
			CallOuter(rt.autoRecover, rt.GetReportError(), rt.startFunc)

			for {
				select {
				case <-startChan:
					invokeStartFun()
					notifyStart()
					rt.GC()

				case callBundle, ok := <-rt.safeCallList:
					if !ok {
						return
					}

					callBundle.Ret <- invokeSafeCallFun(callBundle)
					notifyStart()
					rt.GC()

				case <-rt.Done():
					return
				}
			}

		} else {
			CallOuter(rt.autoRecover, rt.GetReportError(), func() {
				rt.frame = rt.frameCreatorFunc().(FrameWhole)
			})

			var ticker *time.Ticker

			if !rt.frame.IsBlink() {
				ticker = time.NewTicker(time.Duration(float64(time.Second) / float64(rt.frame.GetTargetFPS())))
				defer ticker.Stop()
			}

			uptEntityFun := func() {
				rt.frame.UpdateBegin()
				defer rt.frame.UpdateEnd()

				invokeFun(func(entity EntityWhole) {
					entity.CallUpdate()
				})

				invokeFun(func(entity EntityWhole) {
					entity.CallLateUpdate()
				})

				notifyStart()
				rt.GC()
			}

			loopFun := func() bool {
				if rt.frame.GetTotalFrames() > 0 {
					if rt.frame.GetCurFrames() >= rt.frame.GetTotalFrames() {
						return false
					}
				}

				rt.frame.FrameBegin()
				defer rt.frame.FrameEnd()

				if ticker != nil {
					for {
						select {
						case <-startChan:
							invokeStartFun()
							notifyStart()
							rt.GC()

						case callBundle, ok := <-rt.safeCallList:
							if !ok {
								return false
							}

							callBundle.Ret <- invokeSafeCallFun(callBundle)
							notifyStart()
							rt.GC()

						case <-ticker.C:
							uptEntityFun()
							return true

						case <-rt.Done():
							return false
						}
					}

				} else {
					for {
						select {
						case <-startChan:
							invokeStartFun()
							notifyStart()
							rt.GC()

						case callBundle, ok := <-rt.safeCallList:
							if !ok {
								return false
							}

							callBundle.Ret <- invokeSafeCallFun(callBundle)
							notifyStart()
							rt.GC()

						case <-rt.Done():
							return false

						default:
							uptEntityFun()
							return true
						}
					}
				}
			}

			CallOuter(rt.autoRecover, rt.GetReportError(), rt.startFunc)

			rt.frame.CycleBegin()
			defer rt.frame.CycleEnd()

			for rt.frame.SetCurFrames(0); ; rt.frame.SetCurFrames(rt.frame.GetCurFrames() + 1) {
				if !loopFun() {
					return
				}
			}
		}
	}()

	return rt.shutChan
}

func (rt *Runtime) Stop() {
	rt.GetCancelFunc()()
}

func (rt *Runtime) GetRuntimeID() uint64 {
	return rt.id
}

func (rt *Runtime) GetApp() internal.App {
	return rt.app
}

func (rt *Runtime) GetFrame() internal.Frame {
	return rt.frame
}

func (rt *Runtime) GetEntity(entID uint64) internal.Entity {
	e, _ := rt.entityMap[entID]
	if e.Escape() || e.GetMark(0) {
		return nil
	}

	return e.Value.(internal.Entity)
}

func (rt *Runtime) RangeEntities(fun func(entity internal.Entity) bool) {
	if fun == nil {
		return
	}

	rt.entityList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}
		return fun(e.Value.(internal.Entity))
	})
}

func (rt *Runtime) PushSafeCall(callBundle *SafeCallBundle) {
	if callBundle == nil {
		panic("nil callBundle")
	}

	rt.safeCallList <- callBundle
}

func (rt *Runtime) PushGC(gc internal.GC) {
	if gc != nil {
		rt.gcList = append(rt.gcList, gc)
	}
}

func (rt *Runtime) GC() {
	for i := 0; i < len(rt.entityGCList); i++ {
		rt.entityList.Remove(rt.entityGCList[i])
	}
	rt.entityGCList = rt.entityGCList[:0]

	for i := 0; i < len(rt.gcList); i++ {
		rt.gcList[i].GC()
	}
	rt.gcList = rt.gcList[:0]
}

func (rt *Runtime) AddEntity(entity internal.Entity) {
	if entity == nil {
		panic("nil entity")
	}

	if _, ok := rt.entityMap[entity.GetEntityID()]; ok {
		panic("entity id already exists")
	}

	ele := rt.entityList.PushBack(entity)
	rt.entityMap[entity.GetEntityID()] = ele
	rt.entityStartList = append(rt.entityStartList, ele)
}

func (rt *Runtime) RemoveEntity(entID uint64) {
	if e, ok := rt.entityMap[entID]; ok {
		delete(rt.entityMap, entID)
		e.SetMark(0, true)
		rt.entityGCList = append(rt.entityGCList, e)
	}
}
