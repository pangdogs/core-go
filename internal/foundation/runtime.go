package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
	"time"
	"unsafe"
)

type Runtime interface {
	Runnable
	Context
	GCRoot
	initRuntime(ctx Context, app App, opts *RuntimeOptions)
	GetRuntimeID() uint64
	getInheritor() Runtime
	GetApp() App
	GetFrame() Frame
	GetEntity(entID uint64) Entity
	RangeEntities(fun func(entity Entity) bool)
	GetCache() *misc.Cache
	addEntity(entity Entity)
	removeEntity(entID uint64)
	pushSafeCall(callBundle *SafeCallBundle)
	bindEvent(hookID, eventSrcID uint64, hookEle, eventSrcEle *misc.Element) error
	unbindEvent(hookID, eventSrcID uint64) (hookEle, eventSrcEle *misc.Element, ok bool)
	eventIsBound(hookID, eventSrcID uint64) bool
	eventHandleToBit(handle uintptr) int
	declareEventType(eventID int32, eventType unsafe.Pointer)
	obtainEventType(eventID int32) unsafe.Pointer
}

func RuntimeGetInheritor(rt Runtime) Runtime {
	return rt.getInheritor()
}

func NewRuntime(ctx Context, app App, optFuncs ...NewRuntimeOptionFunc) Runtime {
	opts := &RuntimeOptions{}
	NewRuntimeOption.Default()(opts)

	for _, optFun := range optFuncs {
		optFun(opts)
	}

	var rt *RuntimeFoundation

	if opts.inheritor != nil {
		opts.inheritor.initRuntime(ctx, app, opts)
		return opts.inheritor
	}

	rt = &RuntimeFoundation{}
	rt.initRuntime(ctx, app, opts)

	return rt.inheritor
}

type EventBinderKey struct {
	HookID, EventSrcID uint64
}

type EventBinderValue struct {
	HookEle, EventSrcEle *misc.Element
}

type EventSubscriberKey struct {
	HookID  uint64
	EventID int32
}

type RuntimeFoundation struct {
	RunnableFoundation
	Context
	RuntimeOptions
	id              uint64
	app             App
	safeCallList    chan *SafeCallBundle
	entityMap       map[uint64]*misc.Element
	entityList      misc.List
	entityStartList []*misc.Element
	entityGCList    []*misc.Element
	frame           Frame
	eventBinderMap  map[EventBinderKey]EventBinderValue
	eventHandleBits map[uintptr]int
	eventTypes      [eventsLimit]unsafe.Pointer
	gcExists        map[uintptr]struct{}
	gcList          []GC
	gcLastRunTime   time.Time
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
	rt.id = app.makeUID()
	rt.app = app
	rt.safeCallList = make(chan *SafeCallBundle)
	close(rt.safeCallList)
	rt.entityList.Init(rt.cache)
	rt.entityMap = map[uint64]*misc.Element{}
	rt.eventBinderMap = map[EventBinderKey]EventBinderValue{}
	rt.eventHandleBits = map[uintptr]int{}

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
				CallOuter(rt.autoRecover, rt.GetReportError(), IFace2Entity(e.GetIFace(0)).callStart)
			}
			rt.entityStartList = rt.entityStartList[count:]
		}

		invokeLifecycleFunc := func(fun func(entity Entity)) {
			if fun == nil {
				return
			}

			rt.entityList.UnsafeTraversal(func(e *misc.Element) bool {
				if e.Escape() || e.GetMark(0) {
					return true
				}
				fun(IFace2Entity(e.GetIFace(0)))
				return true
			})
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

		if rt.gcLastRunTime.IsZero() {
			rt.gcLastRunTime = time.Now()
		}
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

				rt.RunGC()
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

				invokeLifecycleFunc(func(entity Entity) {
					entity.callUpdate()
				})

				invokeLifecycleFunc(func(entity Entity) {
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

						rt.RunGC()
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

						rt.RunGC()
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
	if !rt.gcEnable || gc == nil {
		return
	}

	if rt.gcExists == nil {
		rt.gcExists = map[uintptr]struct{}{}
	} else {
		if _, ok := rt.gcExists[gc.GCHandle()]; ok {
			return
		}
	}

	rt.gcExists[gc.GCHandle()] = struct{}{}
	rt.gcList = append(rt.gcList, gc)
}

func (rt *RuntimeFoundation) RunGC() {
	if !rt.gcEnable {
		return
	}

	var gcFlag bool

	if !gcFlag && rt.gcItemNum > 0 {
		if len(rt.gcList) < rt.gcItemNum {
			return
		} else {
			gcFlag = true
		}
	}

	if !gcFlag && rt.gcTimeInterval > 0 {
		if time.Now().Sub(rt.gcLastRunTime) < rt.gcTimeInterval {
			return
		} else {
			gcFlag = true
		}
	}

	if !gcFlag {
		return
	}

	for i := 0; i < len(rt.gcList); i++ {
		rt.gcList[i].GC()
	}

	rt.gcExists = nil
	rt.gcList = rt.gcList[:0]
	rt.gcLastRunTime = time.Now()
}

func (rt *RuntimeFoundation) GC() {
	for i := 0; i < len(rt.entityGCList); i++ {
		rt.entityList.Remove(rt.entityGCList[i])
	}
	rt.entityGCList = rt.entityGCList[:0]
}

func (rt *RuntimeFoundation) GCHandle() uintptr {
	return uintptr(unsafe.Pointer(rt))
}

func (rt *RuntimeFoundation) GCEnabled() bool {
	return rt.gcEnable
}

func (rt *RuntimeFoundation) GetRuntimeID() uint64 {
	return rt.id
}

func (rt *RuntimeFoundation) getInheritor() Runtime {
	return rt.inheritor
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

	return IFace2Entity(e.GetIFace(0))
}

func (rt *RuntimeFoundation) RangeEntities(fun func(entity Entity) bool) {
	if fun == nil {
		return
	}

	rt.entityList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}
		return fun(IFace2Entity(e.GetIFace(0)))
	})
}

func (rt *RuntimeFoundation) GetCache() *misc.Cache {
	return rt.cache
}

func (rt *RuntimeFoundation) addEntity(entity Entity) {
	if entity == nil {
		panic("nil entity")
	}

	if _, ok := rt.entityMap[entity.GetEntityID()]; ok {
		panic("entity id already exists")
	}

	ele := rt.entityList.PushIFaceBack(Entity2IFace(entity))
	rt.entityMap[entity.GetEntityID()] = ele
	rt.entityStartList = append(rt.entityStartList, ele)

}

func (rt *RuntimeFoundation) removeEntity(entID uint64) {
	if e, ok := rt.entityMap[entID]; ok {
		delete(rt.entityMap, entID)
		e.SetMark(0, true)
		if rt.GCEnabled() {
			rt.entityGCList = append(rt.entityGCList, e)
			rt.PushGC(rt)
		}
	}
}

func (rt *RuntimeFoundation) pushSafeCall(callBundle *SafeCallBundle) {
	if callBundle == nil {
		panic("nil callBundle")
	}

	rt.safeCallList <- callBundle
}

func (rt *RuntimeFoundation) bindEvent(hookID, eventSrcID uint64, hookEle, eventSrcEle *misc.Element) error {
	if hookEle == nil {
		return errors.New("nil hookEle")
	}

	if eventSrcEle == nil {
		return errors.New("nil eventSrcEle")
	}

	rt.eventBinderMap[EventBinderKey{
		HookID:     hookID,
		EventSrcID: eventSrcID,
	}] = EventBinderValue{
		HookEle:     hookEle,
		EventSrcEle: eventSrcEle,
	}

	return nil
}

func (rt *RuntimeFoundation) unbindEvent(hookID, eventSrcID uint64) (hookEle, eventSrcEle *misc.Element, ok bool) {
	k := EventBinderKey{
		HookID:     hookID,
		EventSrcID: eventSrcID,
	}

	v, ok := rt.eventBinderMap[k]
	if !ok {
		return nil, nil, false
	}

	delete(rt.eventBinderMap, k)

	return v.HookEle, v.EventSrcEle, true
}

func (rt *RuntimeFoundation) eventIsBound(hookID, eventSrcID uint64) bool {
	_, ok := rt.eventBinderMap[EventBinderKey{
		HookID:     hookID,
		EventSrcID: eventSrcID,
	}]
	return ok
}

func (rt *RuntimeFoundation) eventHandleToBit(handle uintptr) int {
	bit, ok := rt.eventHandleBits[handle]
	if !ok {
		bit = len(rt.eventHandleBits) + 64
		rt.eventHandleBits[handle] = bit
	}
	return bit
}

func (rt *RuntimeFoundation) declareEventType(eventID int32, eventType unsafe.Pointer) {
	if rt.eventTypes[eventID] != nil {
		if rt.eventTypes[eventID] != eventType {
			panic("inconsistent event type")
		}
	}

	rt.eventTypes[eventID] = eventType
}

func (rt *RuntimeFoundation) obtainEventType(eventID int32) unsafe.Pointer {
	if rt.eventTypes[eventID] == nil {
		panic("undeclared event type")
	}

	return rt.eventTypes[eventID]
}
