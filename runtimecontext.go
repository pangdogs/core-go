package core

import "github.com/pangdogs/core/container"

type RuntimeContext interface {
	container.GC
	container.GCCollector
	Context
	_RunnableMark
	EntityMgr
	EntityMgrEvents
	EntityReverseQuery
	SafeCall
	init(appCtx AppContext, opts *RuntimeContextOptions)
	getOptions() *RuntimeContextOptions
	GetAppCtx() AppContext
	setFrame(frame Frame)
	GetFrame() Frame
}

func RuntimeContextGetOptions(runtimeCtx RuntimeContext) RuntimeContextOptions {
	return *runtimeCtx.getOptions()
}

func RuntimeContextGetInheritor(runtimeCtx RuntimeContext) RuntimeContext {
	return runtimeCtx.getOptions().Inheritor
}

func NewRuntimeContext(appCtx AppContext, optFuncs ...NewRuntimeContextOptionFunc) RuntimeContext {
	opts := &RuntimeContextOptions{}
	NewRuntimeContextOption.Default()(opts)

	for i := range optFuncs {
		optFuncs[i](opts)
	}

	var runtimeCtx *RuntimeContextBehavior

	if opts.Inheritor != nil {
		opts.Inheritor.init(appCtx, opts)
		return opts.Inheritor
	}

	runtimeCtx = &RuntimeContextBehavior{}
	runtimeCtx.init(appCtx, opts)

	return runtimeCtx.opts.Inheritor
}

type RuntimeCtxEntityInfo struct {
	Element *container.Element[Face]
	Hooks   [2]Hook
}

type RuntimeContextBehavior struct {
	_ContextBehavior
	_RunnableMarkBehavior
	opts                                RuntimeContextOptions
	appCtx                              AppContext
	entityMap                           map[uint64]RuntimeCtxEntityInfo
	entityList                          container.List[Face]
	frame                               Frame
	callee                              _Callee
	eventEntityMgrAddEntity             Event
	eventEntityMgrRemoveEntity          Event
	eventEntityMgrEntityAddComponents   Event
	eventEntityMgrEntityRemoveComponent Event
	gcMark                              bool
	gcList                              []container.GC
}

func (runtimeCtx *RuntimeContextBehavior) GC() bool {
	if !runtimeCtx.gcMark {
		return false
	}
	runtimeCtx.gcMark = false

	runtimeCtx.entityList.GC()
	runtimeCtx.eventEntityMgrAddEntity.GC()
	runtimeCtx.eventEntityMgrRemoveEntity.GC()
	runtimeCtx.eventEntityMgrEntityAddComponents.GC()
	runtimeCtx.eventEntityMgrEntityRemoveComponent.GC()

	for i := range runtimeCtx.gcList {
		runtimeCtx.gcList[i].GC()
	}
	runtimeCtx.gcList = runtimeCtx.gcList[:0]

	return true
}

func (runtimeCtx *RuntimeContextBehavior) MarkGC() {
	runtimeCtx.gcMark = true
}

func (runtimeCtx *RuntimeContextBehavior) NeedGC() bool {
	return runtimeCtx.gcMark
}

func (runtimeCtx *RuntimeContextBehavior) CollectGC(gc container.GC) {
	if gc == nil || !gc.NeedGC() {
		return
	}

	runtimeCtx.gcList = append(runtimeCtx.gcList, gc)
	runtimeCtx.opts.Inheritor.MarkGC()
}

func (runtimeCtx *RuntimeContextBehavior) init(appCtx AppContext, opts *RuntimeContextOptions) {
	if appCtx == nil {
		panic("nil appCtx")
	}

	if opts == nil {
		panic("nil opts")
	}

	runtimeCtx.opts = *opts

	if runtimeCtx.opts.Inheritor == nil {
		runtimeCtx.opts.Inheritor = runtimeCtx
	}

	runtimeCtx._ContextBehavior.init(appCtx, runtimeCtx.opts.ReportError)
	runtimeCtx.appCtx = appCtx

	runtimeCtx.entityList.Init(runtimeCtx.opts.FaceCache, runtimeCtx.opts.Inheritor)
	runtimeCtx.entityMap = map[uint64]RuntimeCtxEntityInfo{}

	runtimeCtx.eventEntityMgrAddEntity.Init(false, nil, EventRecursion_Discard, runtimeCtx.opts.HookCache, runtimeCtx.opts.Inheritor)
	runtimeCtx.eventEntityMgrRemoveEntity.Init(false, nil, EventRecursion_Discard, runtimeCtx.opts.HookCache, runtimeCtx.opts.Inheritor)
	runtimeCtx.eventEntityMgrEntityAddComponents.Init(false, nil, EventRecursion_Discard, runtimeCtx.opts.HookCache, runtimeCtx.opts.Inheritor)
	runtimeCtx.eventEntityMgrEntityRemoveComponent.Init(false, nil, EventRecursion_Discard, runtimeCtx.opts.HookCache, runtimeCtx.opts.Inheritor)
}

func (runtimeCtx *RuntimeContextBehavior) getOptions() *RuntimeContextOptions {
	return &runtimeCtx.opts
}

func (runtimeCtx *RuntimeContextBehavior) GetAppCtx() AppContext {
	return runtimeCtx.appCtx
}

func (runtimeCtx *RuntimeContextBehavior) setFrame(frame Frame) {
	runtimeCtx.frame = frame
}

func (runtimeCtx *RuntimeContextBehavior) GetFrame() Frame {
	return runtimeCtx.frame
}
