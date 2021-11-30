package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

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

type ComponentHalt interface {
	Halt()
}

type ComponentShut interface {
	Shut(c Component)
}
