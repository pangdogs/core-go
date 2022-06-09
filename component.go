package core

import (
	container2 "github.com/pangdogs/core/container"
)

type Component interface {
	container2.GC
	init(name string, entity Entity, inheritor Component, hookCache *container2.Cache[Hook])
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
}

func (comp *ComponentBehavior) GC() {
	comp.eventComponentDestroySelf.GC()
}

func (comp *ComponentBehavior) init(name string, entity Entity, inheritor Component, hookCache *container2.Cache[Hook]) {
	comp.name = name
	comp.entity = entity
	comp.inheritor = inheritor
	comp.eventComponentDestroySelf.Init(false, nil, hookCache)
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
	EmitEventComponentDestroySelf(&comp.eventComponentDestroySelf, comp.inheritor)
}

func (comp *ComponentBehavior) EventComponentDestroySelf() IEvent {
	return &comp.eventComponentDestroySelf
}
