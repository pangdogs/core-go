package core

import (
	"fmt"
	"github.com/pangdogs/core/container"
)

func (runtimeCtx *RuntimeContextBehavior) GetEntity(id uint64) (Entity, bool) {
	e, ok := runtimeCtx.entityMap[id]
	if !ok {
		return nil, false
	}

	if e.Element.Escaped() {
		return nil, false
	}

	return Fast2IFace[Entity](e.Element.Value.FastIFace), true
}

func (runtimeCtx *RuntimeContextBehavior) RangeEntities(fun func(entity Entity) bool) {
	if fun == nil {
		return
	}

	runtimeCtx.entityList.Traversal(func(e *container.Element[Face]) bool {
		return fun(Fast2IFace[Entity](e.Value.FastIFace))
	})
}

func (runtimeCtx *RuntimeContextBehavior) ReverseRangeEntities(fun func(entity Entity) bool) {
	if fun == nil {
		return
	}

	runtimeCtx.entityList.ReverseTraversal(func(e *container.Element[Face]) bool {
		return fun(Fast2IFace[Entity](e.Value.FastIFace))
	})
}

func (runtimeCtx *RuntimeContextBehavior) AddEntity(entity Entity) {
	if entity == nil {
		panic("nil entity")
	}

	if entity.GetRuntimeCtx() != nil {
		panic("entity already added in runtime context")
	}

	entity.setID(runtimeCtx.appCtx.genUID())
	entity.setRuntimeCtx(runtimeCtx.opts.Inheritor)
	entity.RangeComponents(func(comp Component) bool {
		comp.setID(runtimeCtx.appCtx.genUID())
		return true
	})

	if _, ok := runtimeCtx.entityMap[entity.GetID()]; ok {
		panic(fmt.Errorf("repeated entity '{%d}' in this runtime context", entity.GetID()))
	}

	entityInfo := RuntimeCtxEntityInfo{}

	entityInfo.Hooks[0] = BindEvent[EventCompMgrAddComponents[Entity]](entity.EventCompMgrAddComponents(), runtimeCtx)
	entityInfo.Hooks[1] = BindEvent[EventCompMgrRemoveComponent[Entity]](entity.EventCompMgrRemoveComponent(), runtimeCtx)

	entityInfo.Element = runtimeCtx.entityList.PushBack(Face{
		IFace:     entity,
		FastIFace: IFace2Fast(entity),
	})
	entityInfo.Element.GC = entity

	runtimeCtx.entityMap[entity.GetID()] = entityInfo

	if entity.NeedGC() {
		runtimeCtx.MarkGC()
	}

	emitEventEntityMgrAddEntity[RuntimeContext](&runtimeCtx.eventEntityMgrAddEntity, runtimeCtx.opts.Inheritor, entity)
}

func (runtimeCtx *RuntimeContextBehavior) RemoveEntity(id uint64) {
	e, ok := runtimeCtx.entityMap[id]
	if !ok {
		return
	}

	delete(runtimeCtx.entityMap, id)
	e.Element.Escape()

	for i := range e.Hooks {
		e.Hooks[i].Unbind()
	}

	emitEventEntityMgrRemoveEntity[RuntimeContext](&runtimeCtx.eventEntityMgrRemoveEntity, runtimeCtx.opts.Inheritor, Fast2IFace[Entity](e.Element.Value.FastIFace))
}

func (runtimeCtx *RuntimeContextBehavior) EventEntityMgrAddEntity() IEvent {
	return &runtimeCtx.eventEntityMgrAddEntity
}

func (runtimeCtx *RuntimeContextBehavior) EventEntityMgrRemoveEntity() IEvent {
	return &runtimeCtx.eventEntityMgrRemoveEntity
}

func (runtimeCtx *RuntimeContextBehavior) EventEntityMgrEntityAddComponents() IEvent {
	return &runtimeCtx.eventEntityMgrEntityAddComponents
}

func (runtimeCtx *RuntimeContextBehavior) EventEntityMgrEntityRemoveComponent() IEvent {
	return &runtimeCtx.eventEntityMgrEntityRemoveComponent
}

func (runtimeCtx *RuntimeContextBehavior) OnCompMgrAddComponents(entity Entity, components []Component) {
	for i := range components {
		components[i].setID(runtimeCtx.appCtx.genUID())
	}
	emitEventEntityMgrEntityAddComponents(&runtimeCtx.eventEntityMgrEntityAddComponents, runtimeCtx.opts.Inheritor, entity, components)
}

func (runtimeCtx *RuntimeContextBehavior) OnCompMgrRemoveComponent(entity Entity, component Component) {
	emitEventEntityMgrEntityRemoveComponent(&runtimeCtx.eventEntityMgrEntityRemoveComponent, runtimeCtx.opts.Inheritor, entity, component)
}
