package core

type ComponentAwake interface {
	Awake()
}

type ComponentStart interface {
	Start()
}

type ComponentUpdate = EventUpdate

type ComponentLateUpdate = EventLateUpdate

type ComponentShut interface {
	Shut()
}
