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
	init(servCtx ServiceContext, opts *RuntimeContextOptions)
	getOptions() *RuntimeContextOptions
	GetServiceCtx() ServiceContext
	setFrame(frame Frame)
	GetFrame() Frame
}

func RuntimeContextGetOptions(runtimeCtx RuntimeContext) RuntimeContextOptions {
	return *runtimeCtx.getOptions()
}

func RuntimeContextGetInheritor(runtimeCtx RuntimeContext) RuntimeContext {
	return runtimeCtx.getOptions().Inheritor
}

func NewRuntimeContext(servCtx ServiceContext, optFuncs ...NewRuntimeContextOptionFunc) RuntimeContext {
	opts := &RuntimeContextOptions{}
	NewRuntimeContextOption.Default()(opts)

	for i := range optFuncs {
		optFuncs[i](opts)
	}

	var runtimeCtx *RuntimeContextBehavior

	if opts.Inheritor != nil {
		opts.Inheritor.init(servCtx, opts)
		return opts.Inheritor
	}

	runtimeCtx = &RuntimeContextBehavior{}
	runtimeCtx.init(servCtx, opts)

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
	servCtx                             ServiceContext
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

func (runtimeCtx *RuntimeContextBehavior) init(servCtx ServiceContext, opts *RuntimeContextOptions) {
	if servCtx == nil {
		panic("nil servCtx")
	}

	if opts == nil {
		panic("nil opts")
	}

	runtimeCtx.opts = *opts

	if runtimeCtx.opts.Inheritor == nil {
		runtimeCtx.opts.Inheritor = runtimeCtx
	}

	runtimeCtx._ContextBehavior.init(servCtx, runtimeCtx.opts.ReportError)
	runtimeCtx.servCtx = servCtx

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

func (runtimeCtx *RuntimeContextBehavior) GetServiceCtx() ServiceContext {
	return runtimeCtx.servCtx
}

func (runtimeCtx *RuntimeContextBehavior) setFrame(frame Frame) {
	runtimeCtx.frame = frame
}

func (runtimeCtx *RuntimeContextBehavior) GetFrame() Frame {
	return runtimeCtx.frame
}
