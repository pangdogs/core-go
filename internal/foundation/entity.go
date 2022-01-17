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
	RemoveComponentByID(id uint64)
	GetComponent(name string) Component
	GetComponentByID(id uint64) Component
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
	componentByIDMap        map[uint64]*misc.Element
	componentList           misc.List
	componentGCList         []*misc.Element
	lifecycleEntityInitFunc func()
	lifecycleStartFunc      func()
	lifecycleUpdateFunc     func() bool
	lifecycleLateUpdateFunc func() bool
	lifecycleEntityShutFunc func()
}

func (ent *EntityFoundation) initEntity(rt Runtime, opts *EntityOptions) {
	if rt == nil {
		panic("nil runtime")
	}

	if opts == nil {
		panic("nil opts")
	}

	ent.id = rt.GetApp().makeUID()
	ent.EntityOptions = *opts

	if ent.inheritor == nil {
		ent.inheritor = ent
	}

	ent.runtime = rt
	ent.componentList.Init(rt.GetCache())
	if ent.enableFastGetComponent {
		ent.componentMap = map[string]*misc.Element{}
	}
	if ent.enableFastGetComponentByID {
		ent.componentByIDMap = map[uint64]*misc.Element{}
	}

	rt.GetApp().addEntity(ent.inheritor)
	rt.addEntity(ent.inheritor)

	if ent.initFunc != nil {
		ent.initFunc(ent)
	}

	ent.callEntityInit()
}

func (ent *EntityFoundation) GC() {
	if ent.destroyed {
		return
	}
	for i := range ent.componentGCList {
		ent.componentList.Remove(ent.componentGCList[i])
	}
	ent.componentGCList = ent.componentGCList[:0]
}

func (ent *EntityFoundation) GCHandle() uintptr {
	return uintptr(unsafe.Pointer(ent))
}

func (ent *EntityFoundation) Destroy() {
	if ent.destroyed {
		return
	}

	ent.destroyed = true

	ent.GetRuntime().GetApp().removeEntity(ent.id)
	ent.GetRuntime().removeEntity(ent.id)

	ent.callEntityShut()

	if ent.enableFastGetComponent {
		ent.componentMap = nil
	}

	if ent.enableFastGetComponentByID {
		ent.componentByIDMap = nil
	}

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		e.SetMark(EntityComponentsMark_Removed, true)

		return true
	})

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_HaltedAndShut) {
			return true
		}
		e.SetMark(EntityComponentsMark_HaltedAndShut, true)

		component := IFace2Component(e.GetIFace(EntityComponentsIFace_Component))

		if ch, ok := component.(ComponentHalt); ok {
			ch.Halt()
		}

		if cs, ok := component.(ComponentShut); ok {
			cs.Shut()
		}

		return true
	})

	if ent.shutFunc != nil {
		ent.shutFunc(ent)
	}
}

func (ent *EntityFoundation) GetEntityID() uint64 {
	return ent.id
}

func (ent *EntityFoundation) getEntityInheritor() Entity {
	return ent.inheritor
}

func (ent *EntityFoundation) GetRuntime() Runtime {
	return ent.runtime
}

func (ent *EntityFoundation) IsDestroyed() bool {
	return ent.destroyed
}

func (ent *EntityFoundation) AddComponent(name string, component interface{}) error {
	if component == nil {
		return errors.New("nil component")
	}

	if ent.destroyed {
		return errors.New("entity destroyed")
	}

	_component := component.(Component)
	_component.initComponent(name, ent.inheritor, _component)

	e, ok := ent.getComponentElement(name)
	if ok {
		old := e
		for t := e; t != nil && IFace2Component(t.GetIFace(EntityComponentsIFace_Component)).GetName() == name; t = t.Next() {
			if t.Escape() || t.GetMark(EntityComponentsMark_Removed) {
				continue
			}
			old = t
		}
		e = ent.componentList.InsertIFaceAfter(Component2IFace(_component), old)
	} else {
		e = ent.componentList.PushIFaceBack(Component2IFace(_component))
		if ent.enableFastGetComponent {
			ent.componentMap[name] = e
		}
		if ent.enableFastGetComponentByID {
			ent.componentByIDMap[_component.GetComponentID()] = e
		}
	}

	if ci, ok := _component.(ComponentInit); ok {
		ci.Init()
	}

	if e.Escape() || e.GetMark(EntityComponentsMark_HaltedAndShut) {
		return nil
	}

	if ca, ok := _component.(ComponentAwake); ok {
		ca.Awake()
	}

	return nil
}

func (ent *EntityFoundation) RemoveComponent(name string) {
	e, ok := ent.getComponentElement(name)
	if !ok {
		return
	}

	if ent.enableFastGetComponent {
		delete(ent.componentMap, name)
	}

	for t := e; t != nil; t = t.Next() {
		component := IFace2Component(t.GetIFace(EntityComponentsIFace_Component))
		if component.GetName() != name {
			break
		}

		if t.Escape() || t.GetMark(EntityComponentsMark_Removed) {
			continue
		}
		t.SetMark(EntityComponentsMark_Removed, true)

		if ent.enableFastGetComponentByID {
			delete(ent.componentByIDMap, component.GetComponentID())
		}
	}

	for t := e; t != nil; t = t.Next() {
		component := IFace2Component(t.GetIFace(EntityComponentsIFace_Component))
		if component.GetName() != name {
			break
		}

		if t.Escape() || t.GetMark(EntityComponentsMark_HaltedAndShut) {
			continue
		}
		t.SetMark(EntityComponentsMark_HaltedAndShut, true)

		if ch, ok := component.(ComponentHalt); ok {
			ch.Halt()
		}

		if cs, ok := component.(ComponentShut); ok {
			cs.Shut()
		}

		if !ent.destroyed {
			if ent.runtime.GCEnabled() {
				ent.componentGCList = append(ent.componentGCList, t)
				ent.runtime.PushGC(ent)
			}
		}
	}
}

func (ent *EntityFoundation) RemoveComponentByID(id uint64) {
	e, ok := ent.getComponentElementByID(id)
	if !ok {
		return
	}

	if ent.enableFastGetComponentByID {
		delete(ent.componentByIDMap, id)
	}

	if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
		return
	}
	e.SetMark(EntityComponentsMark_Removed, true)
	e.SetMark(EntityComponentsMark_HaltedAndShut, true)

	component := IFace2Component(e.GetIFace(EntityComponentsIFace_Component))

	if ent.enableFastGetComponent {
		allRemoved := true

		for t := e; t != nil; t = t.Next() {
			other := IFace2Component(t.GetIFace(EntityComponentsIFace_Component))
			if other.GetName() != component.GetName() {
				break
			}

			if t.Escape() || t.GetMark(EntityComponentsMark_Removed) {
				continue
			}

			allRemoved = false
			break
		}

		if allRemoved {
			delete(ent.componentMap, component.GetName())
		}
	}

	if ch, ok := component.(ComponentHalt); ok {
		ch.Halt()
	}

	if cs, ok := component.(ComponentShut); ok {
		cs.Shut()
	}

	if !ent.destroyed {
		if ent.runtime.GCEnabled() {
			ent.componentGCList = append(ent.componentGCList, e)
			ent.runtime.PushGC(ent)
		}
	}
}

func (ent *EntityFoundation) GetComponent(name string) Component {
	if e, ok := ent.getComponentElement(name); ok {
		return IFace2Component(e.GetIFace(EntityComponentsIFace_Component))
	}

	return nil
}

func (ent *EntityFoundation) GetComponentByID(id uint64) Component {
	if e, ok := ent.getComponentElementByID(id); ok {
		return IFace2Component(e.GetIFace(EntityComponentsIFace_Component))
	}

	return nil
}

func (ent *EntityFoundation) GetComponents(name string) []Component {
	if e, ok := ent.getComponentElement(name); ok {
		var components []Component

		for t := e; t != nil && IFace2Component(t.GetIFace(EntityComponentsIFace_Component)).GetName() == name; t = t.Next() {
			if t.Escape() || t.GetMark(EntityComponentsMark_Removed) {
				continue
			}
			components = append(components, IFace2Component(t.GetIFace(EntityComponentsIFace_Component)))
		}

		return components
	}

	return nil
}

func (ent *EntityFoundation) RangeComponents(fun func(component Component) bool) {
	if fun == nil {
		return
	}

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		return fun(IFace2Component(e.GetIFace(EntityComponentsIFace_Component)))
	})
}

func (ent *EntityFoundation) getComponentElement(name string) (*misc.Element, bool) {
	if ent.enableFastGetComponent {
		e, ok := ent.componentMap[name]
		return e, ok
	}

	var element *misc.Element

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).GetName() == name {
			element = e
			return false
		}
		return true
	})

	return element, element != nil
}

func (ent *EntityFoundation) getComponentElementByID(id uint64) (*misc.Element, bool) {
	if ent.enableFastGetComponentByID {
		e, ok := ent.componentByIDMap[id]
		return e, ok
	}

	var element *misc.Element

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).GetComponentID() == id {
			element = e
			return false
		}
		return true
	})

	return element, element != nil
}
