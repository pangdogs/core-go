package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/list"
)

type Entity interface {
	GC
	Destroy()
	GetEntityID() uint64
	GetInheritor() Entity
	GetRuntime() Runtime
	IsDestroyed() bool
	AddComponent(name string, component interface{}) error
	RemoveComponent(name string)
	GetComponent(name string) Component
	GetComponents(name string) []Component
	RangeComponents(fun func(component Component) bool)
}

func NewEntity(rt Runtime, optFuncs ...NewEntityOptionFunc) Entity {
	e := &EntityFoundation{}

	opts := &EntityOptions{}
	NewEntityOption.Default()(opts)

	for _, optFun := range optFuncs {
		optFun(opts)
	}

	e.initEntity(rt, opts)

	return e.inheritor
}

type EntityLifecycleCaller interface {
	CallEntityInit()
	CallStart()
	CallUpdate()
	CallLateUpdate()
	CallEntityShut()
}

const (
	EntityComponentsMark_Removed uint = iota
	EntityComponentsMark_Started
	EntityComponentsMark_NoUpdate
	EntityComponentsMark_NoLateUpdate
)

type EntityFoundation struct {
	EntityOptions
	id              uint64
	runtime         Runtime
	destroyed       bool
	componentMap    map[string]*list.Element
	componentList   list.List
	componentGCList []*list.Element
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
	e.componentList.Init()
	e.componentMap = map[string]*list.Element{}

	rt.GetApp().addEntity(e.inheritor)
	rt.addEntity(e.inheritor)

	if e.initFunc != nil {
		e.initFunc(e)
	}

	e.CallEntityInit()
}

func (e *EntityFoundation) GC() {
	for i := 0; i < len(e.componentGCList); i++ {
		e.componentList.Remove(e.componentGCList[i])
	}
	e.componentGCList = e.componentGCList[:0]
}

func (e *EntityFoundation) Destroy() {
	if e.destroyed {
		return
	}

	e.destroyed = true

	e.GetRuntime().GetApp().removeEntity(e.id)
	e.GetRuntime().removeEntity(e.id)

	e.CallEntityShut()

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

func (e *EntityFoundation) GetInheritor() Entity {
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
	component.initComponent(name, e.inheritor)

	if ele, ok := e.componentMap[name]; ok {
		old := ele
		for t := ele; t != nil && t.Value.(Component).GetName() == name; t = t.Next() {
			old = t
		}
		e.componentList.InsertAfter(_component, old)
	} else {
		e.componentMap[name] = e.componentList.PushBack(_component)
	}

	if cl, ok := _component.(ComponentInit); ok {
		cl.Init(component)
	}

	if cl, ok := _component.(ComponentAwake); ok {
		cl.Awake()
	}

	return nil
}

func (e *EntityFoundation) RemoveComponent(name string) {
	if ele, ok := e.componentMap[name]; ok {
		delete(e.componentMap, name)

		var elements []*list.Element

		for t := ele; t != nil && t.Value.(Component).GetName() == name; t = t.Next() {
			t.SetMark(EntityComponentsMark_Removed, true)
			elements = append(elements, t)
		}

		e.componentGCList = append(e.componentGCList, elements...)

		for i := 0; i < len(elements); i++ {
			if cl, ok := elements[i].Value.(ComponentHalt); ok {
				cl.Halt()
			}
			if cl, ok := elements[i].Value.(ComponentShut); ok {
				cl.Shut(elements[i].Value.(Component))
			}
		}

		if !e.destroyed {
			e.runtime.PushGC(e)
		}
	}
}

func (e *EntityFoundation) GetComponent(name string) Component {
	if ele, ok := e.componentMap[name]; ok {
		return ele.Value.(Component)
	}

	return nil
}

func (e *EntityFoundation) GetComponents(name string) []Component {
	if ele, ok := e.componentMap[name]; ok {
		var components []Component

		for t := ele; t != nil && t.Value.(Component).GetName() == name; t = t.Next() {
			components = append(components, t.Value.(Component))
		}

		return components
	}

	return nil
}

func (e *EntityFoundation) RangeComponents(fun func(component Component) bool) {
	if fun == nil {
		return
	}

	e.componentList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		return fun(e.Value.(Component))
	})
}

func (e *EntityFoundation) CallEntityInit() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if cl, ok := e.Value.(ComponentEntityInit); ok {
			cl.EntityInit()
		}
		return true
	})
}

func (e *EntityFoundation) CallStart() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if !e.GetMark(EntityComponentsMark_Started) {
			e.SetMark(EntityComponentsMark_Started, true)

			if cl, ok := e.Value.(ComponentStart); ok {
				cl.Start()
			}
		}
		return true
	})
}

func (e *EntityFoundation) CallUpdate() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || e.GetMark(EntityComponentsMark_NoUpdate) {
			return true
		}
		if e.GetMark(EntityComponentsMark_Started) {
			if cl, ok := e.Value.(ComponentUpdate); ok {
				cl.Update()
			} else {
				e.SetMark(EntityComponentsMark_NoUpdate, true)
			}
		}
		return true
	})
}

func (e *EntityFoundation) CallLateUpdate() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || e.GetMark(EntityComponentsMark_NoLateUpdate) {
			return true
		}
		if e.GetMark(EntityComponentsMark_Started) {
			if cl, ok := e.Value.(ComponentLateUpdate); ok {
				cl.LateUpdate()
			} else {
				e.SetMark(EntityComponentsMark_NoLateUpdate, true)
			}
		}
		return true
	})
}

func (e *EntityFoundation) CallEntityShut() {
	e.componentList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if cl, ok := e.Value.(ComponentEntityShut); ok {
			cl.EntityShut()
		}
		return true
	})
}
