package internal

type Component interface {
	GetEntity() Entity
	GetName() string
}

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
