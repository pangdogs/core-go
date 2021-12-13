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
	setNotAnalysisLifecycle(v bool)
	getNotAnalysisLifecycle() bool
	setLifecycleComponentInit(ci ComponentInit)
	setLifecycleComponentAwake(ca ComponentAwake)
	setLifecycleComponentEntityInit(cei ComponentEntityInit)
	setLifecycleComponentStart(cs ComponentStart)
	setLifecycleComponentUpdate(cu ComponentUpdate)
	setLifecycleComponentLateUpdate(clu ComponentLateUpdate)
	setLifecycleComponentEntityShut(ces ComponentEntityShut)
	setLifecycleComponentHalt(ch ComponentHalt)
	setLifecycleComponentShut(cs ComponentShut)
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

func ComponentGetInheritor(c Component) Component {
	return c.getComponentInheritor()
}

func ComponentSetNotAnalysisLifecycle(c Component, v bool) {
	c.setNotAnalysisLifecycle(v)
}

func ComponentGetNotAnalysisLifecycle(c Component) bool {
	return c.getNotAnalysisLifecycle()
}

func ComponentSetLifecycleComponentInit(c Component, ci ComponentInit) {
	c.setLifecycleComponentInit(ci)
}

func ComponentSetLifecycleComponentAwake(c Component, ca ComponentAwake) {
	c.setLifecycleComponentAwake(ca)
}

func ComponentSetLifecycleComponentEntityInit(c Component, cei ComponentEntityInit) {
	c.setLifecycleComponentEntityInit(cei)
}

func ComponentSetLifecycleComponentStart(c Component, cs ComponentStart) {
	c.setLifecycleComponentStart(cs)
}

func ComponentSetLifecycleComponentUpdate(c Component, cu ComponentUpdate) {
	c.setLifecycleComponentUpdate(cu)
}

func ComponentSetLifecycleComponentLateUpdate(c Component, clu ComponentLateUpdate) {
	c.setLifecycleComponentLateUpdate(clu)
}

func ComponentSetLifecycleComponentEntityShut(c Component, ces ComponentEntityShut) {
	c.setLifecycleComponentEntityShut(ces)
}

func ComponentSetLifecycleComponentHalt(c Component, ch ComponentHalt) {
	c.setLifecycleComponentHalt(ch)
}

func ComponentSetLifecycleComponentShut(c Component, cs ComponentShut) {
	c.setLifecycleComponentShut(cs)
}

func ComponentGetLifecycleComponentInit(c Component) ComponentInit {
	return c.getLifecycleComponentInit()
}

func ComponentGetLifecycleComponentAwake(c Component) ComponentAwake {
	return c.getLifecycleComponentAwake()
}

func ComponentGetLifecycleComponentEntityInit(c Component) ComponentEntityInit {
	return c.getLifecycleComponentEntityInit()
}

func ComponentGetLifecycleComponentStart(c Component) ComponentStart {
	return c.getLifecycleComponentStart()
}

func ComponentGetLifecycleComponentUpdate(c Component) ComponentUpdate {
	return c.getLifecycleComponentUpdate()
}

func ComponentGetLifecycleComponentLateUpdate(c Component) ComponentLateUpdate {
	return c.getLifecycleComponentLateUpdate()
}

func ComponentGetLifecycleComponentEntityShut(c Component) ComponentEntityShut {
	return c.getLifecycleComponentEntityShut()
}

func ComponentGetLifecycleComponentHalt(c Component) ComponentHalt {
	return c.getLifecycleComponentHalt()
}

func ComponentGetLifecycleComponentShut(c Component) ComponentShut {
	return c.getLifecycleComponentShut()
}

type ComponentFoundation struct {
	name                         string
	entity                       Entity
	inheritor                    Component
	notAnalysisLifecycle         bool
	lifecycleComponentInit       ComponentInit
	lifecycleComponentAwake      ComponentAwake
	lifecycleComponentEntityInit ComponentEntityInit
	lifecycleComponentStart      ComponentStart
	lifecycleComponentUpdate     ComponentUpdate
	lifecycleComponentLateUpdate ComponentLateUpdate
	lifecycleComponentEntityShut ComponentEntityShut
	lifecycleComponentHalt       ComponentHalt
	lifecycleComponentShut       ComponentShut
}

func (c *ComponentFoundation) initComponent(name string, entity Entity, inheritor Component) {
	c.name = name
	c.entity = entity
	c.inheritor = inheritor

	if !c.notAnalysisLifecycle {
		if ci, ok := inheritor.(ComponentInit); ok {
			c.setLifecycleComponentInit(ci)
		}
		if ca, ok := inheritor.(ComponentAwake); ok {
			c.setLifecycleComponentAwake(ca)
		}
		if cei, ok := inheritor.(ComponentEntityInit); ok {
			c.setLifecycleComponentEntityInit(cei)
		}
		if cs, ok := inheritor.(ComponentStart); ok {
			c.setLifecycleComponentStart(cs)
		}
		if cu, ok := inheritor.(ComponentUpdate); ok {
			c.setLifecycleComponentUpdate(cu)
		}
		if clu, ok := inheritor.(ComponentLateUpdate); ok {
			c.setLifecycleComponentLateUpdate(clu)
		}
		if ces, ok := inheritor.(ComponentEntityShut); ok {
			c.setLifecycleComponentEntityShut(ces)
		}
		if ch, ok := inheritor.(ComponentHalt); ok {
			c.setLifecycleComponentHalt(ch)
		}
		if cs, ok := inheritor.(ComponentShut); ok {
			c.setLifecycleComponentShut(cs)
		}
	}
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
