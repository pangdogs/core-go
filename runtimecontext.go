package core

import (
	container2 "github.com/pangdogs/core/container"
)

type RuntimeContext interface {
	container2.GC
	Context
	EntityMgr
	EntityMgrEvents
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

	for _, optFunc := range optFuncs {
		optFunc(opts)
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
	Element *container2.Element[Face]
	Hooks   [2]Hook
}

type RuntimeContextBehavior struct {
	ContextBehavior
	opts                                RuntimeContextOptions
	appCtx                              AppContext
	entityMap                           map[uint64]RuntimeCtxEntityInfo
	entityList                          container2.List[Face]
	frame                               Frame
	eventEntityMgrAddEntity             Event
	eventEntityMgrRemoveEntity          Event
	eventEntityMgrEntityAddComponents   Event
	eventEntityMgrEntityRemoveComponent Event
	eventPushSafeCallSegment            Event
}

func (runtimeCtx *RuntimeContextBehavior) GC() {
	runtimeCtx.entityList.GC()
	runtimeCtx.eventEntityMgrAddEntity.GC()
	runtimeCtx.eventEntityMgrRemoveEntity.GC()
	runtimeCtx.eventEntityMgrEntityAddComponents.GC()
	runtimeCtx.eventEntityMgrEntityRemoveComponent.GC()
	runtimeCtx.eventPushSafeCallSegment.GC()
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

	runtimeCtx.ContextBehavior.init(appCtx)
	runtimeCtx.appCtx = appCtx

	runtimeCtx.entityList.Init(runtimeCtx.opts.FaceCache)
	runtimeCtx.entityMap = map[uint64]RuntimeCtxEntityInfo{}

	runtimeCtx.eventEntityMgrAddEntity.Init(false, nil, runtimeCtx.opts.HookCache)
	runtimeCtx.eventEntityMgrRemoveEntity.Init(false, nil, runtimeCtx.opts.HookCache)
	runtimeCtx.eventEntityMgrEntityAddComponents.Init(false, nil, runtimeCtx.opts.HookCache)
	runtimeCtx.eventEntityMgrEntityRemoveComponent.Init(false, nil, runtimeCtx.opts.HookCache)
	runtimeCtx.eventPushSafeCallSegment.Init(false, nil, runtimeCtx.opts.HookCache)
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
