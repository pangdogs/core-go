package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type Component interface {
	GetName() string
	GetEntity() Entity
	GetComponentInheritor() Component
	initComponent(name string, entity Entity, inheritor Component)
	getLifecycleComponentInit() ComponentInit
	getLifecycleComponentAwake() ComponentAwake
	getLifecycleComponentEntityInit() ComponentEntityInit
	getLifecycleComponentStart() ComponentStart
	getLifecycleComponentUpdate() ComponentUpdate
	getLifecycleComponentLateUpdate() ComponentLateUpdate
	getLifecycleComponentEntityShut() ComponentEntityShut
	getLifecycleComponentHalt() ComponentHalt
	getLifecycleComponentShut() ComponentShut
}

func IFace2Component(f misc.IFace) Component {
	return *(*Component)(unsafe.Pointer(&f))
}

func Component2IFace(c Component) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&c))
}

type ComponentFoundation struct {
	name                 string
	entity               Entity
	inheritor            Component
	lifecycleTab         [ComponentLifecycle_Count]misc.IFace
	NotAnalysisLifecycle bool
}

func (c *ComponentFoundation) GetName() string {
	return c.name
}

func (c *ComponentFoundation) GetEntity() Entity {
	return c.entity
}

func (c *ComponentFoundation) GetComponentInheritor() Component {
	return c.inheritor
}

func (c *ComponentFoundation) initComponent(name string, entity Entity, inheritor Component) {
	c.name = name
	c.entity = entity
	c.inheritor = inheritor

	if !c.NotAnalysisLifecycle {
		if ci, ok := inheritor.(ComponentInit); ok {
			c.SetLifecycleComponentInit(ci)
		}
		if ca, ok := inheritor.(ComponentAwake); ok {
			c.SetLifecycleComponentAwake(ca)
		}
		if cei, ok := inheritor.(ComponentEntityInit); ok {
			c.SetLifecycleComponentEntityInit(cei)
		}
		if cs, ok := inheritor.(ComponentStart); ok {
			c.SetLifecycleComponentStart(cs)
		}
		if cu, ok := inheritor.(ComponentUpdate); ok {
			c.SetLifecycleComponentUpdate(cu)
		}
		if clu, ok := inheritor.(ComponentLateUpdate); ok {
			c.SetLifecycleComponentLateUpdate(clu)
		}
		if ces, ok := inheritor.(ComponentEntityShut); ok {
			c.SetLifecycleComponentEntityShut(ces)
		}
		if ch, ok := inheritor.(ComponentHalt); ok {
			c.SetLifecycleComponentHalt(ch)
		}
		if cs, ok := inheritor.(ComponentShut); ok {
			c.SetLifecycleComponentShut(cs)
		}
	}
}
