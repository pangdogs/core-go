package internal

type App interface {
	Runnable
	Context
	GetEntity(entID uint64) Entity
}
