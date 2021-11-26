package foundation

import (
	"unsafe"
)

type Component interface {
	GetEntity() Entity
	GetName() string
	initComponent(name string, entity Entity)
}

func IFace2Component(p unsafe.Pointer) Component {
	return *(*Component)(p)
}

func Component2IFace(c Component) unsafe.Pointer {
	return unsafe.Pointer(&c)
}

type ComponentFoundation struct {
	name   string
	entity Entity
}

func (c *ComponentFoundation) GetEntity() Entity {
	return c.entity
}

func (c *ComponentFoundation) GetName() string {
	return c.name
}

func (c *ComponentFoundation) initComponent(name string, entity Entity) {
	c.name = name
	c.entity = entity
}
