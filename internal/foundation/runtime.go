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
	eventRecursionEnabled() bool
	recursiveEventDiscarded() bool
	incrEventCalledDepth() bool
	decrEventCalledDepth()
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

type RuntimeFoundation struct {
	RunnableFoundation
	Context
	RuntimeOptions
	id               uint64
	app              App
	safeCallList     chan *SafeCallBundle
	entityMap        map[uint64]*misc.Element
	entityList       misc.List
	entityStartList  []*misc.Element
	entityGCList     []*misc.Element
	frame            Frame
	eventBinderMap   map[EventBinderKey]EventBinderValue
	eventCalledDepth int
	gcExists         map[uintptr]struct{}
	gcList           []GC
	gcLastRunTime    time.Time
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

	rt.Context = ctx
	rt.id = app.makeUID()
	rt.app = app
	rt.safeCallList = make(chan *SafeCallBundle)
	close(rt.safeCallList)
	rt.entityList.Init(rt.cache)
	rt.entityMap = map[uint64]*misc.Element{}
	rt.eventBinderMap = map[EventBinderKey]EventBinderValue{}

	CallOuter(rt.enableAutoRecover, rt.GetReportError(), func() {
		if rt.initFunc != nil {
			rt.initFunc(rt)
		}
	})

	if opts.enableAutoRun {
		rt.Run()
	}
}

func (rt *RuntimeFoundation) Run() chan struct{} {
	if !rt.markRunning() {
		panic("runtime already running")
	}

	rt.safeCallList = make(chan *SafeCallBundle, rt.safeCallCacheSize)
	shutChan := make(chan struct{}, 1)

	go rt.running(shutChan)

	return shutChan
}

func (rt *RuntimeFoundation) Stop() {
	rt.GetCancelFunc()()
}

func (rt *RuntimeFoundation) PushGC(gc GC) {
	if !rt.enableGC || gc == nil {
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
	if !rt.enableGC {
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

	for i := range rt.gcList {
		rt.gcList[i].GC()
	}

	rt.gcExists = nil
	rt.gcList = rt.gcList[:0]
	if rt.gcTimeInterval > 0 {
		rt.gcLastRunTime = time.Now()
	}
}

func (rt *RuntimeFoundation) GC() {
	for i := range rt.entityGCList {
		rt.entityList.Remove(rt.entityGCList[i])
	}
	rt.entityGCList = rt.entityGCList[:0]
}

func (rt *RuntimeFoundation) GCHandle() uintptr {
	return uintptr(unsafe.Pointer(rt))
}

func (rt *RuntimeFoundation) GCEnabled() bool {
	return rt.enableGC
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

	e := rt.entityList.PushIFaceBack(Entity2IFace(entity))
	rt.entityMap[entity.GetEntityID()] = e
	rt.entityStartList = append(rt.entityStartList, e)
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

func (rt *RuntimeFoundation) eventRecursionEnabled() bool {
	return rt.enableEventRecursion
}

func (rt *RuntimeFoundation) recursiveEventDiscarded() bool {
	return rt.discardRecursiveEvent
}

func (rt *RuntimeFoundation) incrEventCalledDepth() bool {
	if rt.callEventDepth > 0 {
		if rt.eventCalledDepth > rt.callEventDepth {
			if rt.discardExceedDepthEvent {
				return false
			}
			panic("event called exceed limited depth")
		}
	}

	rt.eventCalledDepth++

	return true
}

func (rt *RuntimeFoundation) decrEventCalledDepth() {
	if rt.eventCalledDepth > 0 {
		rt.eventCalledDepth--
	}
}
