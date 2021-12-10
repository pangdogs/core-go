package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type ComponentInit interface {
	Init()
}

func IFace2ComponentInit(f misc.IFace) ComponentInit {
	return *(*ComponentInit)(unsafe.Pointer(&f))
}

func ComponentInit2IFace(ci ComponentInit) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&ci))
}

type ComponentAwake interface {
	Awake()
}

func IFace2ComponentAwake(f misc.IFace) ComponentAwake {
	return *(*ComponentAwake)(unsafe.Pointer(&f))
}

func ComponentAwake2IFace(ca ComponentAwake) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&ca))
}

type ComponentEntityInit interface {
	EntityInit()
}

func IFace2ComponentEntityInit(f misc.IFace) ComponentEntityInit {
	return *(*ComponentEntityInit)(unsafe.Pointer(&f))
}

func ComponentEntityInit2IFace(cei ComponentEntityInit) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&cei))
}

type ComponentStart interface {
	Start()
}

func IFace2ComponentStart(f misc.IFace) ComponentStart {
	return *(*ComponentStart)(unsafe.Pointer(&f))
}

func ComponentStart2IFace(cs ComponentStart) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&cs))
}

type ComponentUpdate interface {
	Update()
}

func IFace2ComponentUpdate(f misc.IFace) ComponentUpdate {
	return *(*ComponentUpdate)(unsafe.Pointer(&f))
}

func ComponentUpdate2IFace(cu ComponentUpdate) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&cu))
}

type ComponentLateUpdate interface {
	LateUpdate()
}

func IFace2ComponentLateUpdate(f misc.IFace) ComponentLateUpdate {
	return *(*ComponentLateUpdate)(unsafe.Pointer(&f))
}

func ComponentLateUpdate2IFace(clu ComponentLateUpdate) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&clu))
}

type ComponentEntityShut interface {
	EntityShut()
}

func IFace2ComponentEntityShut(f misc.IFace) ComponentEntityShut {
	return *(*ComponentEntityShut)(unsafe.Pointer(&f))
}

func ComponentEntityShut2IFace(ces ComponentEntityShut) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&ces))
}

type ComponentHalt interface {
	Halt()
}

func IFace2ComponentHalt(f misc.IFace) ComponentHalt {
	return *(*ComponentHalt)(unsafe.Pointer(&f))
}

func ComponentHalt2IFace(ch ComponentHalt) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&ch))
}

type ComponentShut interface {
	Shut()
}

func IFace2ComponentShut(f misc.IFace) ComponentShut {
	return *(*ComponentShut)(unsafe.Pointer(&f))
}

func ComponentShut2IFace(cs ComponentShut) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&cs))
}

type ComponentLifecycle int

const (
	ComponentLifecycle_ComponentInit ComponentLifecycle = iota
	ComponentLifecycle_ComponentAwake
	ComponentLifecycle_ComponentEntityInit
	ComponentLifecycle_ComponentStart
	ComponentLifecycle_ComponentUpdate
	ComponentLifecycle_ComponentLateUpdate
	ComponentLifecycle_ComponentEntityShut
	ComponentLifecycle_ComponentHalt
	ComponentLifecycle_ComponentShut
	ComponentLifecycle_Count
)

func (c *ComponentFoundation) setNotAnalysisLifecycle(v bool) {
	c.notAnalysisLifecycle = v
}

func (c *ComponentFoundation) getNotAnalysisLifecycle() bool {
	return c.notAnalysisLifecycle
}

func (c *ComponentFoundation) setLifecycleIFace(lifecycle ComponentLifecycle, face misc.IFace) {
	c.lifecycleTab[lifecycle] = face
}

func (c *ComponentFoundation) getLifecycleIFace(lifecycle ComponentLifecycle) misc.IFace {
	return c.lifecycleTab[lifecycle]
}

func (c *ComponentFoundation) setLifecycleComponentInit(ci ComponentInit) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentInit, ComponentInit2IFace(ci))
}

func (c *ComponentFoundation) setLifecycleComponentAwake(ca ComponentAwake) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentAwake, ComponentAwake2IFace(ca))
}

func (c *ComponentFoundation) setLifecycleComponentEntityInit(cei ComponentEntityInit) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentEntityInit, ComponentEntityInit2IFace(cei))
}

func (c *ComponentFoundation) setLifecycleComponentStart(cs ComponentStart) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentStart, ComponentStart2IFace(cs))
}

func (c *ComponentFoundation) setLifecycleComponentUpdate(cu ComponentUpdate) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentUpdate, ComponentUpdate2IFace(cu))
}

func (c *ComponentFoundation) setLifecycleComponentLateUpdate(clu ComponentLateUpdate) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentLateUpdate, ComponentLateUpdate2IFace(clu))
}

func (c *ComponentFoundation) setLifecycleComponentEntityShut(ces ComponentEntityShut) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentEntityShut, ComponentEntityShut2IFace(ces))
}

func (c *ComponentFoundation) setLifecycleComponentHalt(ch ComponentHalt) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentHalt, ComponentHalt2IFace(ch))
}

func (c *ComponentFoundation) setLifecycleComponentShut(cs ComponentShut) {
	c.setLifecycleIFace(ComponentLifecycle_ComponentShut, ComponentShut2IFace(cs))
}

func (c *ComponentFoundation) getLifecycleComponentInit() ComponentInit {
	return IFace2ComponentInit(c.getLifecycleIFace(ComponentLifecycle_ComponentInit))
}

func (c *ComponentFoundation) getLifecycleComponentAwake() ComponentAwake {
	return IFace2ComponentAwake(c.getLifecycleIFace(ComponentLifecycle_ComponentAwake))
}

func (c *ComponentFoundation) getLifecycleComponentEntityInit() ComponentEntityInit {
	return IFace2ComponentEntityInit(c.getLifecycleIFace(ComponentLifecycle_ComponentEntityInit))
}

func (c *ComponentFoundation) getLifecycleComponentStart() ComponentStart {
	return IFace2ComponentStart(c.getLifecycleIFace(ComponentLifecycle_ComponentStart))
}

func (c *ComponentFoundation) getLifecycleComponentUpdate() ComponentUpdate {
	return IFace2ComponentUpdate(c.getLifecycleIFace(ComponentLifecycle_ComponentUpdate))
}

func (c *ComponentFoundation) getLifecycleComponentLateUpdate() ComponentLateUpdate {
	return IFace2ComponentLateUpdate(c.getLifecycleIFace(ComponentLifecycle_ComponentLateUpdate))
}

func (c *ComponentFoundation) getLifecycleComponentEntityShut() ComponentEntityShut {
	return IFace2ComponentEntityShut(c.getLifecycleIFace(ComponentLifecycle_ComponentEntityShut))
}

func (c *ComponentFoundation) getLifecycleComponentHalt() ComponentHalt {
	return IFace2ComponentHalt(c.getLifecycleIFace(ComponentLifecycle_ComponentHalt))
}

func (c *ComponentFoundation) getLifecycleComponentShut() ComponentShut {
	return IFace2ComponentShut(c.getLifecycleIFace(ComponentLifecycle_ComponentShut))
}
