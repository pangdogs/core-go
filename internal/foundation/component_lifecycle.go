package foundation

type ComponentInit interface {
	Init(c Component)
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
	Shut(c Component)
}
