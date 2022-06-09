package core

import (
	"unsafe"
)

type FastIFace [2]unsafe.Pointer

func Fast2IFace[T any](fi FastIFace) T {
	return *(*T)(unsafe.Pointer(&fi))
}

func IFace2Fast[T any](i T) FastIFace {
	return *(*FastIFace)(unsafe.Pointer(&i))
}

var NilFastIFace FastIFace

type Face struct {
	IFace     interface{}
	FastIFace FastIFace
}
