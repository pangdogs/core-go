package internal

type App interface {
	Runnable
	Context
	GetEntity(entID uint64) Entity
	RangeEntities(func(entity Entity) bool)
	MakeUID() uint64
}
