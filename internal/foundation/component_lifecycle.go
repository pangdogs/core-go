package foundation

import "unsafe"

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

func IFace2ComponentUpdate(p unsafe.Pointer) ComponentUpdate {
	return *(*ComponentUpdate)(p)
}

func ComponentUpdate2IFace(cu ComponentUpdate) unsafe.Pointer {
	return unsafe.Pointer(&cu)
}

type ComponentLateUpdate interface {
	LateUpdate()
}

func IFace2ComponentLateUpdate(p unsafe.Pointer) ComponentLateUpdate {
	return *(*ComponentLateUpdate)(p)
}

func ComponentLateUpdate2IFace(clu ComponentLateUpdate) unsafe.Pointer {
	return unsafe.Pointer(&clu)
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
