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
	setPrimary(v bool)
	getPrimary() bool
	DestroySelf()
	eventComponentDestroySelf() IEvent
}

func ComponentGetInheritor(comp Component) Component {
	return comp.getInheritor()
}

type ComponentBehavior struct {
	id                         uint64
	name                       string
	entity                     Entity
	inheritor                  Component
	primary                    bool
	_eventComponentDestroySelf Event
	gcMark                     bool
}

func (comp *ComponentBehavior) GC() bool {
	if !comp.gcMark {
		return false
	}
	comp.gcMark = false

	comp._eventComponentDestroySelf.GC()

	return true
}

func (comp *ComponentBehavior) MarkGC() {
	if comp.gcMark {
		return
	}
	comp.gcMark = true

	if comp.entity != nil {
		comp.entity.MarkGC()
	}
}

func (comp *ComponentBehavior) NeedGC() bool {
	return comp.gcMark
}

func (comp *ComponentBehavior) CollectGC(gc container.GC) {
	if gc == nil || !gc.NeedGC() {
		return
	}

	comp.inheritor.MarkGC()
}

func (comp *ComponentBehavior) init(name string, entity Entity, inheritor Component, hookCache *container.Cache[Hook]) {
	comp.name = name
	comp.entity = entity
	comp.inheritor = inheritor
	comp._eventComponentDestroySelf.Init(false, nil, EventRecursion_Discard, hookCache, comp)
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

func (comp *ComponentBehavior) setPrimary(v bool) {
	comp.primary = v
}

func (comp *ComponentBehavior) getPrimary() bool {
	return comp.primary
}

func (comp *ComponentBehavior) DestroySelf() {
	emitEventComponentDestroySelf(&comp._eventComponentDestroySelf, comp.inheritor)
}

func (comp *ComponentBehavior) eventComponentDestroySelf() IEvent {
	return &comp._eventComponentDestroySelf
}
