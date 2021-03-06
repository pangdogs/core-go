package core

import "github.com/pangdogs/core/container"

type Entity interface {
	container.GC
	ComponentMgr
	ComponentMgrEvents
	init(opts *EntityOptions)
	getOptions() *EntityOptions
	setID(id uint64)
	GetID() uint64
	setRuntimeCtx(runtimeCtx RuntimeContext)
	GetRuntimeCtx() RuntimeContext
	DestroySelf()
	eventEntityDestroySelf() IEvent
}

func EntityGetOptions(e Entity) EntityOptions {
	return *e.getOptions()
}

func EntityGetInheritor(e Entity) Entity {
	return e.getOptions().Inheritor
}

func NewEntity(optFuncs ...NewEntityOptionFunc) Entity {
	opts := &EntityOptions{}
	NewEntityOption.Default()(opts)

	for i := range optFuncs {
		optFuncs[i](opts)
	}

	if opts.Inheritor != nil {
		opts.Inheritor.init(opts)
		return opts.Inheritor
	}

	e := &EntityBehavior{}
	e.init(opts)

	return e.opts.Inheritor
}

type EntityBehavior struct {
	id                          uint64
	opts                        EntityOptions
	runtimeCtx                  RuntimeContext
	componentList               container.List[Face]
	componentMap                map[string]*container.Element[Face]
	componentByIDMap            map[uint64]*container.Element[Face]
	_eventEntityDestroySelf     Event
	eventCompMgrAddComponents   Event
	eventCompMgrRemoveComponent Event
	gcMark                      bool
}

func (entity *EntityBehavior) GC() bool {
	if !entity.gcMark {
		return false
	}
	entity.gcMark = false

	entity.componentList.GC()
	entity._eventEntityDestroySelf.GC()
	entity.eventCompMgrAddComponents.GC()
	entity.eventCompMgrRemoveComponent.GC()

	return true
}

func (entity *EntityBehavior) MarkGC() {
	if entity.gcMark {
		return
	}
	entity.gcMark = true

	if entity.runtimeCtx != nil {
		entity.runtimeCtx.MarkGC()
	}
}

func (entity *EntityBehavior) NeedGC() bool {
	return entity.gcMark
}

func (entity *EntityBehavior) CollectGC(gc container.GC) {
	if gc == nil || !gc.NeedGC() {
		return
	}

	entity.opts.Inheritor.MarkGC()
}

func (entity *EntityBehavior) init(opts *EntityOptions) {
	if opts == nil {
		panic("nil opts")
	}

	entity.opts = *opts

	if entity.opts.Inheritor == nil {
		entity.opts.Inheritor = entity
	}

	entity.componentList.Init(entity.opts.FaceCache, entity)

	if entity.opts.EnableFastGetComponent {
		entity.componentMap = map[string]*container.Element[Face]{}
	}

	if entity.opts.EnableFastGetComponentByID {
		entity.componentByIDMap = map[uint64]*container.Element[Face]{}
	}

	entity._eventEntityDestroySelf.Init(false, nil, EventRecursion_Discard, opts.HookCache, entity)
	entity.eventCompMgrAddComponents.Init(false, nil, EventRecursion_Discard, opts.HookCache, entity)
	entity.eventCompMgrRemoveComponent.Init(false, nil, EventRecursion_Discard, opts.HookCache, entity)
}

func (entity *EntityBehavior) getOptions() *EntityOptions {
	return &entity.opts
}

func (entity *EntityBehavior) setID(id uint64) {
	entity.id = id
}

func (entity *EntityBehavior) GetID() uint64 {
	return entity.id
}

func (entity *EntityBehavior) setRuntimeCtx(runtimeCtx RuntimeContext) {
	entity.runtimeCtx = runtimeCtx
}

func (entity *EntityBehavior) GetRuntimeCtx() RuntimeContext {
	return entity.runtimeCtx
}

func (entity *EntityBehavior) DestroySelf() {
	emitEventEntityDestroySelf(&entity._eventEntityDestroySelf, entity.opts.Inheritor)
}

func (entity *EntityBehavior) eventEntityDestroySelf() IEvent {
	return &entity._eventEntityDestroySelf
}
