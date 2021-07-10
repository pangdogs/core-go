package internal

type Runtime interface {
	Runnable
	Context
	GCRoot
	GetRuntimeID() uint64
	GetApp() App
	GetFrame() Frame
	GetEntity(entID uint64) Entity
	RangeEntities(fun func(entity Entity) bool)
}
