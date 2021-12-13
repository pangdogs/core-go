package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type ComponentInit interface {
	Init()
}

func IFace2ComponentInit(f misc.IFace) ComponentInit {
	return *(*ComponentInit)(unsafe.Pointer(&f))
}

func ComponentInit2IFace(ci ComponentInit) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&ci))
}

type ComponentAwake interface {
	Awake()
}

func IFace2ComponentAwake(f misc.IFace) ComponentAwake {
	return *(*ComponentAwake)(unsafe.Pointer(&f))
}

func ComponentAwake2IFace(ca ComponentAwake) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&ca))
}

type ComponentEntityInit interface {
	EntityInit()
}

func IFace2ComponentEntityInit(f misc.IFace) ComponentEntityInit {
	return *(*ComponentEntityInit)(unsafe.Pointer(&f))
}

func ComponentEntityInit2IFace(cei ComponentEntityInit) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&cei))
}

type ComponentStart interface {
	Start()
}

func IFace2ComponentStart(f misc.IFace) ComponentStart {
	return *(*ComponentStart)(unsafe.Pointer(&f))
}

func ComponentStart2IFace(cs ComponentStart) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&cs))
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

func IFace2ComponentEntityShut(f misc.IFace) ComponentEntityShut {
	return *(*ComponentEntityShut)(unsafe.Pointer(&f))
}

func ComponentEntityShut2IFace(ces ComponentEntityShut) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&ces))
}

type ComponentHalt interface {
	Halt()
}

func IFace2ComponentHalt(f misc.IFace) ComponentHalt {
	return *(*ComponentHalt)(unsafe.Pointer(&f))
}

func ComponentHalt2IFace(ch ComponentHalt) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&ch))
}

type ComponentShut interface {
	Shut()
}

func IFace2ComponentShut(f misc.IFace) ComponentShut {
	return *(*ComponentShut)(unsafe.Pointer(&f))
}

func ComponentShut2IFace(cs ComponentShut) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&cs))
}
