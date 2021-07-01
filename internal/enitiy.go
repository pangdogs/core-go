package internal

type Entity interface {
	Destroy()
	GetEntityID() uint64
	GetRuntime() Runtime
	IsDestroyed() bool
	AddComponent(name string, component interface{}) error
	RemoveComponent(name string)
	GetComponent(name string) Component
	GetComponents(name string) []Component
	RangeComponents(fun func(component Component) bool)
}
