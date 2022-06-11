package core

import "github.com/pangdogs/core/container"

type Component interface {
	container.GC
	init(name string, entity Entity, inheritor Component, hookCache *container.Cache[Hook])
	setID(id uint64)
	GetID() uint64
	GetName() string
	GetEntity() Entity
	getInheritor() Component
	DestroySelf()
	EventComponentDestroySelf() IEvent
}

func ComponentGetInheritor(comp Component) Component {
	return comp.getInheritor()
}

type ComponentBehavior struct {
	id                        uint64
	name                      string
	entity                    Entity
	inheritor                 Component
	eventComponentDestroySelf Event
	gcMark                    bool
}

func (comp *ComponentBehavior) GC() {
	if !comp.gcMark {
		return
	}
	comp.gcMark = false

	comp.eventComponentDestroySelf.GC()
}

func (comp *ComponentBehavior) MarkGC() {
	comp.gcMark = true
}

func (comp *ComponentBehavior) NeedGC() bool {
	return comp.gcMark
}

func (comp *ComponentBehavior) init(name string, entity Entity, inheritor Component, hookCache *container.Cache[Hook]) {
	comp.name = name
	comp.entity = entity
	comp.inheritor = inheritor
	comp.eventComponentDestroySelf.Init(false, nil, hookCache, comp)
}

func (comp *ComponentBehavior) setID(id uint64) {
	comp.id = id
}

func (comp *ComponentBehavior) GetID() uint64 {
	return comp.id
}

func (comp *ComponentBehavior) GetName() string {
	return comp.name
}

func (comp *ComponentBehavior) GetEntity() Entity {
	return comp.entity
}

func (comp *ComponentBehavior) getInheritor() Component {
	return comp.inheritor
}

func (comp *ComponentBehavior) DestroySelf() {
	emitEventComponentDestroySelf(&comp.eventComponentDestroySelf, comp.inheritor)
}

func (comp *ComponentBehavior) EventComponentDestroySelf() IEvent {
	return &comp.eventComponentDestroySelf
}
