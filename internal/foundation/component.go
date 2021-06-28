package foundation

import "github.com/pangdogs/core/internal"

type ComponentWhole interface {
	internal.Component
	initComponent(name string, entity internal.Entity)
	setStarted(v bool)
	getStarted() bool
}

type Component struct {
	name    string
	entity  internal.Entity
	started bool
}

func (c *Component) GetEntity() internal.Entity {
	return c.entity
}

func (c *Component) GetName() string {
	return c.name
}

func (c *Component) initComponent(name string, entity internal.Entity) {
	c.name = name
	c.entity = entity
}

func (c *Component) setStarted(v bool) {
	c.started = v
}

func (c *Component) getStarted() bool {
	return c.started
}
