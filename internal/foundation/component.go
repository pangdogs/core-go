package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type Component interface {
	GetEntity() Entity
	GetName() string
	initComponent(name string, entity Entity)
}

func IFace2Component(f misc.IFace) Component {
	return *(*Component)(unsafe.Pointer(&f))
}

func Component2IFace(c Component) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&c))
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
