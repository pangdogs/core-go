package internal

type Component interface {
	GetEntity() Entity
	GetName() string
}

type ComponentAwake interface {
	Awake()
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

type ComponentDestroy interface {
	Destroy()
}
