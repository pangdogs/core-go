package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/list"
	"unsafe"
)

type Entity interface {
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
	getEntityFoundation() *EntityFoundation
}

func NewEntity(rt Runtime, optFuncs ...NewEntityOptionFunc) Entity {
	opts := &EntityOptions{}
	NewEntityOption.Default()(opts)

	for _, optFun := range optFuncs {
		optFun(opts)
	}

	var e *EntityFoundation

	if opts.inheritor != nil {
		e = opts.inheritor.getEntityFoundation()
	} else {
		e = &EntityFoundation{}
	}

	e.initEntity(rt, opts)

	return e.inheritor
}

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
			if t.Escape() || t.GetMark(EntityComponentsMark_Removed) {
				continue
			}
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

		for i := 0; i < len(elements); i++ {
			if cl, ok := elements[i].Value.(ComponentHalt); ok {
				cl.Halt()
			}
			if cl, ok := elements[i].Value.(ComponentShut); ok {
				cl.Shut(elements[i].Value.(Component))
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
	if ele, ok := e.componentMap[name]; ok {
		return ele.Value.(Component)
	}

	return nil
}

func (e *EntityFoundation) GetComponents(name string) []Component {
	if ele, ok := e.componentMap[name]; ok {
		var components []Component

		for t := ele; t != nil && t.Value.(Component).GetName() == name; t = t.Next() {
			if t.Escape() || t.GetMark(EntityComponentsMark_Removed) {
				continue
			}
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

func (e *EntityFoundation) getEntityFoundation() *EntityFoundation {
	return e
}
