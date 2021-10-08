package foundation

type Component interface {
	GetEntity() Entity
	GetName() string
	initComponent(name string, entity Entity)
}

type ComponentFoundation = _Component

type _Component struct {
	name   string
	entity Entity
}

func (c *_Component) GetEntity() Entity {
	return c.entity
}

func (c *_Component) GetName() string {
	return c.name
}

func (c *_Component) initComponent(name string, entity Entity) {
	c.name = name
	c.entity = entity
}
