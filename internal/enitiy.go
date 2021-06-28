package internal

type Entity interface {
	Destroy()
	GetEntityID() uint64
	GetRuntime() Runtime
	IsDestroyed() bool
	AddComponent(name string, component Component) error
	RemoveComponent(name string)
	GetComponent(name string) Component
	GetComponents(name string) []Component
	RangeComponents(fun func(component Component) bool)
}
