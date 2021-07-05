package internal

type Runtime interface {
	Runnable
	Context
	GCRoot
	GetRuntimeID() uint64
	GetApp() App
	GetFrame() Frame
}
