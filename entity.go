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
	EventEntityDestroySelf() IEvent
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

	for _, optFunc := range optFuncs {
		optFunc(opts)
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
	eventEntityDestroySelf      Event
	eventCompMgrAddComponents   Event
	eventCompMgrRemoveComponent Event
}

func (entity *EntityBehavior) GC() {
	entity.componentList.GC()
	entity.eventEntityDestroySelf.GC()
	entity.eventCompMgrAddComponents.GC()
	entity.eventCompMgrRemoveComponent.GC()
}

func (entity *EntityBehavior) init(opts *EntityOptions) {
	if opts == nil {
		panic("nil opts")
	}

	entity.opts = *opts

	if entity.opts.Inheritor == nil {
		entity.opts.Inheritor = entity
	}

	entity.componentList.Init(entity.opts.FaceCache)

	if entity.opts.EnableFastGetComponent {
		entity.componentMap = map[string]*container.Element[Face]{}
	}

	if entity.opts.EnableFastGetComponentByID {
		entity.componentByIDMap = map[uint64]*container.Element[Face]{}
	}

	entity.eventEntityDestroySelf.Init(false, nil, opts.HookCache)
	entity.eventCompMgrAddComponents.Init(false, nil, opts.HookCache)
	entity.eventCompMgrRemoveComponent.Init(false, nil, opts.HookCache)
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
	emitEventEntityDestroySelf(&entity.eventEntityDestroySelf, entity.opts.Inheritor)
}

func (entity *EntityBehavior) EventEntityDestroySelf() IEvent {
	return &entity.eventEntityDestroySelf
}
