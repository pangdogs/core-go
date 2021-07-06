package internal

type Component interface {
	GetEntity() Entity
	GetName() string
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

type ComponentDestroy interface {
	Destroy()
}
