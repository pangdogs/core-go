package foundation

type ComponentInit interface {
	Init()
}

type ComponentAwake interface {
	Awake()
}

type ComponentEntityInit interface {
	EntityInit()
}

type ComponentStart interface {
	Start()
}

type ComponentUpdate interface {
	Update()
}

type ComponentLateUpdate interface {
	LateUpdate()
}

type ComponentEntityShut interface {
	EntityShut()
}

type ComponentHalt interface {
	Halt()
}

type ComponentShut interface {
	Shut()
}

func (c *ComponentFoundation) setNotAnalysisLifecycle(v bool) {
	c.notAnalysisLifecycle = v
}

func (c *ComponentFoundation) getNotAnalysisLifecycle() bool {
	return c.notAnalysisLifecycle
}

func (c *ComponentFoundation) setLifecycleComponentInit(ci ComponentInit) {
	c.lifecycleComponentInit = ci
}

func (c *ComponentFoundation) setLifecycleComponentAwake(ca ComponentAwake) {
	c.lifecycleComponentAwake = ca
}

func (c *ComponentFoundation) setLifecycleComponentEntityInit(cei ComponentEntityInit) {
	c.lifecycleComponentEntityInit = cei
}

func (c *ComponentFoundation) setLifecycleComponentStart(cs ComponentStart) {
	c.lifecycleComponentStart = cs
}

func (c *ComponentFoundation) setLifecycleComponentUpdate(cu ComponentUpdate) {
	c.lifecycleComponentUpdate = cu
}

func (c *ComponentFoundation) setLifecycleComponentLateUpdate(clu ComponentLateUpdate) {
	c.lifecycleComponentLateUpdate = clu
}

func (c *ComponentFoundation) setLifecycleComponentEntityShut(ces ComponentEntityShut) {
	c.lifecycleComponentEntityShut = ces
}

func (c *ComponentFoundation) setLifecycleComponentHalt(ch ComponentHalt) {
	c.lifecycleComponentHalt = ch
}

func (c *ComponentFoundation) setLifecycleComponentShut(cs ComponentShut) {
	c.lifecycleComponentShut = cs
}

func (c *ComponentFoundation) getLifecycleComponentInit() ComponentInit {
	return c.lifecycleComponentInit
}

func (c *ComponentFoundation) getLifecycleComponentAwake() ComponentAwake {
	return c.lifecycleComponentAwake
}

func (c *ComponentFoundation) getLifecycleComponentEntityInit() ComponentEntityInit {
	return c.lifecycleComponentEntityInit
}

func (c *ComponentFoundation) getLifecycleComponentStart() ComponentStart {
	return c.lifecycleComponentStart
}

func (c *ComponentFoundation) getLifecycleComponentUpdate() ComponentUpdate {
	return c.lifecycleComponentUpdate
}

func (c *ComponentFoundation) getLifecycleComponentLateUpdate() ComponentLateUpdate {
	return c.lifecycleComponentLateUpdate
}

func (c *ComponentFoundation) getLifecycleComponentEntityShut() ComponentEntityShut {
	return c.lifecycleComponentEntityShut
}

func (c *ComponentFoundation) getLifecycleComponentHalt() ComponentHalt {
	return c.lifecycleComponentHalt
}

func (c *ComponentFoundation) getLifecycleComponentShut() ComponentShut {
	return c.lifecycleComponentShut
}
