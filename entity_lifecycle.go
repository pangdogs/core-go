package core

type EntityInit interface {
	Init()
}

type EntityInitFin interface {
	InitFin()
}

type EntityUpdate = EventUpdate

type EntityLateUpdate = EventLateUpdate

type EntityShut interface {
	Shut()
}

type EntityShutFin interface {
	ShutFin()
}
