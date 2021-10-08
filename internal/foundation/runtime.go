package foundation

import (
	"github.com/pangdogs/core/internal/list"
	"time"
)

type Runtime interface {
	Runnable
	Context
	GCRoot
	GC
	GetRuntimeID() uint64
	GetApp() App
	GetFrame() Frame
	GetEntity(entID uint64) Entity
	RangeEntities(fun func(entity Entity) bool)
	GetInheritor() Runtime
	addEntity(entity Entity)
	removeEntity(entID uint64)
	pushSafeCall(callBundle *SafeCallBundle)
}

func NewRuntime(ctx Context, app App, optFuncs ...NewRuntimeOptionFunc) Runtime {
	rt := &RuntimeFoundation{}

	opts := &RuntimeOptions{}
	NewRuntimeOption.Default()(opts)

	for _, optFun := range optFuncs {
		optFun(opts)
	}

	rt.initRuntime(ctx, app, opts)

	return rt.inheritor
}

type RuntimeFoundation struct {
	_Runnable
	Context
	RuntimeOptions
	id              uint64
	app             App
	safeCallList    chan *SafeCallBundle
	entityMap       map[uint64]*list.Element
	entityList      list.List
	entityStartList []*list.Element
	entityGCList    []*list.Element
	gcList          []GC
	frame           Frame
}

func (rt *RuntimeFoundation) initRuntime(ctx Context, app App, opts *RuntimeOptions) {
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

	if rt.inheritor == nil {
		rt.inheritor = rt
	}

	rt.initRunnable()
	rt.Context = ctx
	rt.id = app.MakeUID()
	rt.app = app
	rt.safeCallList = make(chan *SafeCallBundle)
	close(rt.safeCallList)
	rt.entityList.Init()
	rt.entityMap = map[uint64]*list.Element{}

	CallOuter(rt.autoRecover, rt.GetReportError(), func() {
		if rt.initFunc != nil {
			rt.initFunc(rt)
		}
	})

	if opts.autoRun {
		rt.Run()
	}
}

func (rt *RuntimeFoundation) Run() chan struct{} {
	if !rt.markRunning() {
		panic("runtime already running")
	}

	rt.safeCallList = make(chan *SafeCallBundle, rt.safeCallCacheSize)

	go func() {
		if parentCtx, ok := rt.GetParentContext().(Context); ok {
			parentCtx.GetWaitGroup().Add(1)
		}

		invokeFun := func(fun func(entity Entity)) {
			if fun == nil {
				return
			}
			rt.RangeEntities(func(entity Entity) bool {
				CallOuter(rt.autoRecover, rt.GetReportError(), func() {
					fun(entity)
				})
				return true
			})
		}

		startChan := make(chan struct{}, 1)

		notifyStart := func() {
			if len(rt.entityStartList) > 0 {
				select {
				case startChan <- struct{}{}:
				default:
				}
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
				CallOuter(rt.autoRecover, rt.GetReportError(), e.Value.(Entity).callStart)
			}
			rt.entityStartList = rt.entityStartList[count:]
		}

		invokeSafeCallFun := func(callBundle *SafeCallBundle) (ret SafeRet) {
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
					notifyStart()

					select {
					case <-startChan:
						invokeStartFun()

					case callBundle, ok := <-rt.safeCallList:
						if !ok {
							return
						}
						callBundle.Ret <- invokeSafeCallFun(callBundle)

					default:
						return
					}

					rt.GC()
				}
			}()

			if parentCtx, ok := rt.GetParentContext().(Context); ok {
				parentCtx.GetWaitGroup().Done()
			}

			rt.GetWaitGroup().Wait()
			rt.markShutdown()
			rt.shutChan <- struct{}{}

			CallOuter(rt.autoRecover, rt.GetReportError(), func() {
				if rt.stopFunc != nil {
					rt.stopFunc(rt)
				}
			})
		}()

		rt.frame = nil

		if rt.frameCreatorFunc == nil {
			CallOuter(rt.autoRecover, rt.GetReportError(), func() {
				if rt.startFunc != nil {
					rt.startFunc(rt)
				}
			})

			for {
				notifyStart()

				select {
				case <-startChan:
					invokeStartFun()

				case callBundle, ok := <-rt.safeCallList:
					if !ok {
						return
					}
					callBundle.Ret <- invokeSafeCallFun(callBundle)

				case <-rt.Done():
					return
				}

				rt.GC()
			}

		} else {
			CallOuter(rt.autoRecover, rt.GetReportError(), func() {
				rt.frame = rt.frameCreatorFunc(rt)
			})

			var ticker *time.Ticker

			if !rt.frame.IsBlink() {
				ticker = time.NewTicker(time.Duration(float64(time.Second) / float64(rt.frame.GetTargetFPS())))
				defer ticker.Stop()
			}

			uptEntityFun := func() {
				rt.frame.updateBegin()
				defer rt.frame.updateEnd()

				invokeFun(func(entity Entity) {
					entity.callUpdate()
				})

				invokeFun(func(entity Entity) {
					entity.callLateUpdate()
				})
			}

			loopFun := func() bool {
				if rt.frame.GetTotalFrames() > 0 {
					if rt.frame.GetCurFrames() >= rt.frame.GetTotalFrames() {
						return false
					}
				}

				rt.frame.frameBegin()
				defer rt.frame.frameEnd()

				if ticker != nil {
					onceUpdate := false

					for {
						notifyStart()

						select {
						case <-startChan:
							invokeStartFun()

						case callBundle, ok := <-rt.safeCallList:
							if !ok {
								return false
							}
							callBundle.Ret <- invokeSafeCallFun(callBundle)

						case <-ticker.C:
							if onceUpdate {
								return true
							}
							onceUpdate = true

							uptEntityFun()

						case <-rt.Done():
							return false
						}

						rt.GC()
					}

				} else {
					onceUpdate := false

					for {
						notifyStart()

						select {
						case <-startChan:
							invokeStartFun()

						case callBundle, ok := <-rt.safeCallList:
							if !ok {
								return false
							}
							callBundle.Ret <- invokeSafeCallFun(callBundle)

						case <-rt.Done():
							return false

						default:
							if onceUpdate {
								return true
							}
							onceUpdate = true

							uptEntityFun()
						}

						rt.GC()
					}
				}
			}

			CallOuter(rt.autoRecover, rt.GetReportError(), func() {
				if rt.startFunc != nil {
					rt.startFunc(rt)
				}
			})

			rt.frame.cycleBegin()
			defer rt.frame.cycleEnd()

			for ; ; rt.frame.setCurFrames(rt.frame.GetCurFrames() + 1) {
				if !loopFun() {
					return
				}
			}
		}
	}()

	return rt.shutChan
}

func (rt *RuntimeFoundation) Stop() {
	rt.GetCancelFunc()()
}

func (rt *RuntimeFoundation) PushGC(gc GC) {
	if gc != nil {
		rt.gcList = append(rt.gcList, gc)
	}
}

func (rt *RuntimeFoundation) GC() {
	for i := 0; i < len(rt.entityGCList); i++ {
		rt.entityList.Remove(rt.entityGCList[i])
	}
	rt.entityGCList = rt.entityGCList[:0]

	for i := 0; i < len(rt.gcList); i++ {
		rt.gcList[i].GC()
	}
	rt.gcList = rt.gcList[:0]
}

func (rt *RuntimeFoundation) GetRuntimeID() uint64 {
	return rt.id
}

func (rt *RuntimeFoundation) GetApp() App {
	return rt.app
}

func (rt *RuntimeFoundation) GetFrame() Frame {
	return rt.frame
}

func (rt *RuntimeFoundation) GetEntity(entID uint64) Entity {
	e, ok := rt.entityMap[entID]
	if !ok {
		return nil
	}

	if e.Escape() || e.GetMark(0) {
		return nil
	}

	return e.Value.(Entity)
}

func (rt *RuntimeFoundation) RangeEntities(fun func(entity Entity) bool) {
	if fun == nil {
		return
	}

	rt.entityList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}
		return fun(e.Value.(Entity))
	})
}

func (rt *RuntimeFoundation) GetInheritor() Runtime {
	return rt.inheritor
}

func (rt *RuntimeFoundation) addEntity(entity Entity) {
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

func (rt *RuntimeFoundation) removeEntity(entID uint64) {
	if e, ok := rt.entityMap[entID]; ok {
		delete(rt.entityMap, entID)
		e.SetMark(0, true)
		rt.entityGCList = append(rt.entityGCList, e)
	}
}

func (rt *RuntimeFoundation) pushSafeCall(callBundle *SafeCallBundle) {
	if callBundle == nil {
		panic("nil callBundle")
	}

	rt.safeCallList <- callBundle
}
