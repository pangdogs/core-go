package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type Component interface {
	initComponent(name string, entity Entity, inheritor Component)
	GetName() string
	GetEntity() Entity
	getComponentInheritor() Component
}

func IFace2Component(f misc.IFace) Component {
	return *(*Component)(unsafe.Pointer(&f))
}

func Component2IFace(c Component) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&c))
}

func ComponentGetInheritor(c Component) Component {
	return c.getComponentInheritor()
}

type ComponentFoundation struct {
	name      string
	entity    Entity
	inheritor Component
}

func (c *ComponentFoundation) initComponent(name string, entity Entity, inheritor Component) {
	c.name = name
	c.entity = entity
	c.inheritor = inheritor
}

func (c *ComponentFoundation) GetName() string {
	return c.name
}

func (c *ComponentFoundation) GetEntity() Entity {
	return c.entity
}

func (c *ComponentFoundation) getComponentInheritor() Component {
	return c.inheritor
}
