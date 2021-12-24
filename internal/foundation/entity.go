package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type Entity interface {
	initEntity(rt Runtime, opts *EntityOptions)
	Destroy()
	GetEntityID() uint64
	getEntityInheritor() Entity
	GetRuntime() Runtime
	IsDestroyed() bool
	AddComponent(name string, component interface{}) error
	RemoveComponent(name string)
	GetComponent(name string) Component
	GetComponents(name string) []Component
	RangeComponents(fun func(component Component) bool)
	callEntityInit()
	callStart()
	callUpdate()
	callLateUpdate()
	callEntityShut()
	setLifecycleEntityInitFunc(fun func())
	setLifecycleStartFunc(fun func())
	setLifecycleUpdateFunc(fun func() bool)
	setLifecycleLateUpdateFunc(fun func() bool)
	setLifecycleEntityShutFunc(fun func())
	getLifecycleEntityInitFunc() func()
	getLifecycleStartFunc() func()
	getLifecycleUpdateFunc() func() bool
	getLifecycleLateUpdateFunc() func() bool
	getLifecycleEntityShutFunc() func()
}

func IFace2Entity(f misc.IFace) Entity {
	return *(*Entity)(unsafe.Pointer(&f))
}

func Entity2IFace(e Entity) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&e))
}

func EntityGetInheritor(e Entity) Entity {
	return e.getEntityInheritor()
}

func EntitySetLifecycleEntityInitFunc(e Entity, fun func()) {
	e.setLifecycleEntityInitFunc(fun)
}

func EntitySetLifecycleStartFunc(e Entity, fun func()) {
	e.setLifecycleStartFunc(fun)
}

func EntitySetLifecycleUpdateFunc(e Entity, fun func() bool) {
	e.setLifecycleUpdateFunc(fun)
}

func EntitySetLifecycleLateUpdateFunc(e Entity, fun func() bool) {
	e.setLifecycleLateUpdateFunc(fun)
}

func EntitySetLifecycleEntityShutFunc(e Entity, fun func()) {
	e.setLifecycleEntityShutFunc(fun)
}

func EntityGetLifecycleEntityInitFunc(e Entity) func() {
	return e.getLifecycleEntityInitFunc()
}

func EntityGetLifecycleStartFunc(e Entity) func() {
	return e.getLifecycleStartFunc()
}

func EntityGetLifecycleUpdateFunc(e Entity) func() bool {
	return e.getLifecycleUpdateFunc()
}

func EntityGetLifecycleLateUpdateFunc(e Entity) func() bool {
	return e.getLifecycleLateUpdateFunc()
}

func EntityGetLifecycleEntityShutFunc(e Entity) func() {
	return e.getLifecycleEntityShutFunc()
}

func NewEntity(rt Runtime, optFuncs ...NewEntityOptionFunc) Entity {
	opts := &EntityOptions{}
	NewEntityOption.Default()(opts)

	for _, optFun := range optFuncs {
		optFun(opts)
	}

	if opts.inheritor != nil {
		opts.inheritor.initEntity(rt, opts)
		return opts.inheritor
	}

	e := &EntityFoundation{}
	e.initEntity(rt, opts)

	return e.inheritor
}

const (
	EntityComponentsIFace_Component int = iota
	EntityComponentsIFace_ComponentUpdate
	EntityComponentsIFace_ComponentLateUpdate
)

type EntityFoundation struct {
	EntityOptions
	id                      uint64
	runtime                 Runtime
	destroyed               bool
	componentMap            map[string]*misc.Element
	componentList           misc.List
	componentGCList         []*misc.Element
	lifecycleEntityInitFunc func()
	lifecycleStartFunc      func()
	lifecycleUpdateFunc     func() bool
	lifecycleLateUpdateFunc func() bool
	lifecycleEntityShutFunc func()
}

func (e *EntityFoundation) initEntity(rt Runtime, opts *EntityOptions) {
	if rt == nil {
		panic("nil runtime")
	}

	if opts == nil {
		panic("nil opts")
	}

	e.id = rt.GetApp().makeUID()
	e.EntityOptions = *opts

	if e.inheritor == nil {
		e.inheritor = e
	}

	e.runtime = rt
	e.componentList.Init(rt.GetCache())
	if e.enableFastGetComponent {
		e.componentMap = map[string]*misc.Element{}
	}

	rt.GetApp().addEntity(e.inheritor)
	rt.addEntity(e.inheritor)

	if e.initFunc != nil {
		e.initFunc(e)
	}

	e.callEntityInit()
}

func (e *EntityFoundation) GC() {
	for i := 0; i < len(e.componentGCList); i++ {
		e.componentList.Remove(e.componentGCList[i])
	}
	e.componentGCList = e.componentGCList[:0]
}

func (e *EntityFoundation) GCHandle() uintptr {
	return uintptr(unsafe.Pointer(e))
}

func (e *EntityFoundation) Destroy() {
	if e.destroyed {
		return
	}

	e.destroyed = true

	e.GetRuntime().GetApp().removeEntity(e.id)
	e.GetRuntime().removeEntity(e.id)

	e.callEntityShut()

	e.RangeComponents(func(component Component) bool {
		e.RemoveComponent(component.GetName())
		return true
	})

	if e.shutFunc != nil {
		e.shutFunc(e)
	}
}

func (e *EntityFoundation) GetEntityID() uint64 {
	return e.id
}

func (e *EntityFoundation) getEntityInheritor() Entity {
	return e.inheritor
}

func (e *EntityFoundation) GetRuntime() Runtime {
	return e.runtime
}

func (e *EntityFoundation) IsDestroyed() bool {
	return e.destroyed
}

func (e *EntityFoundation) AddComponent(name string, _component interface{}) error {
	if _component == nil {
		return errors.New("nil component")
	}

	if e.destroyed {
		return errors.New("entity destroyed")
	}

	component := _component.(Component)
	component.initComponent(name, e.inheritor, component)

	if ele, ok := e.getComponentElement(name); ok {
		old := ele
		for t := ele; t != nil && IFace2Component(t.GetIFace(EntityComponentsIFace_Component)).GetName() == name; t = t.Next() {
			if t.Escape() || t.GetMark(EntityComponentsMark_Removed) {
				continue
			}
			old = t
		}
		e.componentList.InsertIFaceAfter(Component2IFace(component), old)
	} else {
		ele = e.componentList.PushIFaceBack(Component2IFace(component))
		if e.enableFastGetComponent {
			e.componentMap[name] = ele
		}
	}

	if ci, ok := component.(ComponentInit); ok {
		ci.Init()
	}

	if ca, ok := component.(ComponentAwake); ok {
		ca.Awake()
	}

	return nil
}

func (e *EntityFoundation) RemoveComponent(name string) {
	if ele, ok := e.getComponentElement(name); ok {
		if e.enableFastGetComponent {
			delete(e.componentMap, name)
		}

		var elements []*misc.Element

		for t := ele; t != nil && IFace2Component(t.GetIFace(EntityComponentsIFace_Component)).GetName() == name; t = t.Next() {
			t.SetMark(EntityComponentsMark_Removed, true)
			elements = append(elements, t)
		}

		for i := 0; i < len(elements); i++ {
			c := IFace2Component(elements[i].GetIFace(EntityComponentsIFace_Component))

			if ch, ok := c.(ComponentHalt); ok {
				ch.Halt()
			}

			if cs, ok := c.(ComponentShut); ok {
				cs.Shut()
			}
		}

		if !e.destroyed {
			if e.runtime.GCEnabled() {
				e.componentGCList = append(e.componentGCList, elements...)
				e.runtime.PushGC(e)
			}
		}
	}
}

func (e *EntityFoundation) GetComponent(name string) Component {
	if ele, ok := e.getComponentElement(name); ok {
		return IFace2Component(ele.GetIFace(EntityComponentsIFace_Component))
	}

	return nil
}

func (e *EntityFoundation) GetComponents(name string) []Component {
	if ele, ok := e.getComponentElement(name); ok {
		var components []Component

		for t := ele; t != nil && IFace2Component(t.GetIFace(EntityComponentsIFace_Component)).GetName() == name; t = t.Next() {
			if t.Escape() || t.GetMark(EntityComponentsMark_Removed) {
				continue
			}
			components = append(components, IFace2Component(t.GetIFace(EntityComponentsIFace_Component)))
		}

		return components
	}

	return nil
}

func (e *EntityFoundation) RangeComponents(fun func(component Component) bool) {
	if fun == nil {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		return fun(IFace2Component(e.GetIFace(EntityComponentsIFace_Component)))
	})
}

func (e *EntityFoundation) getComponentElement(name string) (*misc.Element, bool) {
	if e.enableFastGetComponent {
		ele, ok := e.componentMap[name]
		return ele, ok
	}

	var ele *misc.Element

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).GetName() == name {
			ele = e
			return false
		}
		return true
	})

	return ele, ele != nil
}
