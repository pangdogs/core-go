package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type Component interface {
	initComponent(name string, entity Entity, inheritor Component)
	shutComponent()
	IsEmbedded() bool
	GetComponentID() uint64
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
	id        uint64
	name      string
	entity    Entity
	inheritor Component
	embedded  bool
}

func (c *ComponentFoundation) initComponent(name string, entity Entity, inheritor Component) {
	if c.entity != nil {
		panic("init repeated")
	}

	c.name = name
	c.entity = entity
	c.inheritor = inheritor
	c.id = entity.GetRuntime().GetApp().makeUID()
	c.embedded = true
}

func (c *ComponentFoundation) shutComponent() {
	c.embedded = false
}

func (c *ComponentFoundation) IsEmbedded() bool {
	return c.embedded
}

func (c *ComponentFoundation) GetComponentID() uint64 {
	return c.id
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
