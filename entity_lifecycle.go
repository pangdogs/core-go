package core

type EntityInit interface {
	Init()
}

type EntityUpdate = EventUpdate

type EntityLateUpdate = EventLateUpdate

type EntityShut interface {
	Shut()
}
