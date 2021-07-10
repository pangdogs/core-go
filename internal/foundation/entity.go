package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal"
	"github.com/pangdogs/core/internal/list"
)

type EntityWhole interface {
	internal.Entity
	internal.GC
	GetInheritor() internal.Entity
	CallStart()
	CallUpdate()
	CallLateUpdate()
}

func NewEntity(rt internal.Runtime, optFuncs ...NewEntityOptionFunc) internal.Entity {
	e := &Entity{}

	opts := &EntityOptions{}
	for _, optFun := range append([]NewEntityOptionFunc{NewEntityOption.Default()}, optFuncs...) {
		optFun(opts)
	}

	e.InitEntity(rt, opts)

	return e.inheritor
}

type Entity struct {
	EntityOptions
	id              uint64
	runtime         internal.Runtime
	destroyed       bool
	componentMap    map[string]*list.Element
	componentList   list.List
	componentGCList []*list.Element
}

func (e *Entity) InitEntity(rt internal.Runtime, opts *EntityOptions) {
	if rt == nil {
		panic("nil runtime")
	}

	if opts == nil {
		panic("nil opts")
	}

	e.id = rt.GetApp().MakeUID()
	e.EntityOptions = *opts

	if e.inheritor != nil {
		e.inheritor.(EntityInheritorWhole).initEntityInheritor(e)
	} else {
		e.inheritor = e
	}

	e.runtime = rt
	e.componentList.Init()
	e.componentMap = map[string]*list.Element{}

	rt.GetApp().(AppWhole).AddEntity(e.inheritor)
	rt.(RuntimeWhole).AddEntity(e.inheritor)

	if e.initFunc != nil {
		e.initFunc()
	}

	e.RangeComponents(func(component internal.Component) bool {
		if cl, ok := component.(internal.ComponentEntityInit); ok {
			cl.EntityInit()
		}
		return true
	})
}

func (e *Entity) Destroy() {
	if e.destroyed {
		return
	}

	e.destroyed = true

	e.GetRuntime().GetApp().(AppWhole).RemoveEntity(e.id)
	e.GetRuntime().(RuntimeWhole).RemoveEntity(e.id)

	e.RangeComponents(func(component internal.Component) bool {
		if cl, ok := component.(internal.ComponentEntityShut); ok {
			cl.EntityShut()
		}
		return true
	})

	e.RangeComponents(func(component internal.Component) bool {
		e.RemoveComponent(component.GetName())
		return true
	})

	if e.shutFunc != nil {
		e.shutFunc()
	}
}

func (e *Entity) GetEntityID() uint64 {
	return e.id
}

func (e *Entity) GetRuntime() internal.Runtime {
	return e.runtime
}

func (e *Entity) IsDestroyed() bool {
	return e.destroyed
}

func (e *Entity) AddComponent(name string, component interface{}) error {
	if name == "" {
		return errors.New("empty component name")
	}

	if component == nil {
		return errors.New("nil component")
	}

	if e.destroyed {
		return errors.New("entity destroyed")
	}

	component.(ComponentWhole).initComponent(name, e.inheritor)

	if ele, ok := e.componentMap[name]; ok {
		old := ele
		for t := ele; t != nil && t.Value.(internal.Component).GetName() == name; t = t.Next() {
			old = t
		}
		e.componentList.InsertAfter(component, old)
	} else {
		e.componentMap[name] = e.componentList.PushBack(component)
	}

	if cl, ok := component.(internal.ComponentAwake); ok {
		cl.Awake()
	}

	return nil
}

func (e *Entity) RemoveComponent(name string) {
	if ele, ok := e.componentMap[name]; ok {
		delete(e.componentMap, name)

		var elements []*list.Element
		for t := ele; t != nil && t.Value.(internal.Component).GetName() == name; t = t.Next() {
			t.SetMark(0, true)
			elements = append(elements, t)
		}

		e.componentGCList = append(e.componentGCList, elements...)

		for i := 0; i < len(elements); i++ {
			if cl, ok := elements[i].Value.(internal.ComponentShut); ok {
				cl.Shut()
			}
		}

		if !e.destroyed {
			e.runtime.PushGC(e)
		}
	}
}

func (e *Entity) GetComponent(name string) internal.Component {
	if ele, ok := e.componentMap[name]; ok {
		return ele.Value.(internal.Component)
	}

	return nil
}

func (e *Entity) GetComponents(name string) []internal.Component {
	if ele, ok := e.componentMap[name]; ok {
		var components []internal.Component

		for t := ele; t != nil && t.Value.(internal.Component).GetName() == name; t = t.Next() {
			components = append(components, t.Value.(internal.Component))
		}

		return components
	}

	return nil
}

func (e *Entity) RangeComponents(fun func(component internal.Component) bool) {
	if fun == nil {
		return
	}

	e.componentList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}
		return fun(e.Value.(internal.Component))
	})
}

func (e *Entity) GC() {
	for i := 0; i < len(e.componentGCList); i++ {
		e.componentList.Remove(e.componentGCList[i])
	}
	e.componentGCList = e.componentGCList[:0]
}

func (e *Entity) GetInheritor() internal.Entity {
	return e.inheritor
}

func (e *Entity) CallStart() {
	if e.destroyed {
		return
	}

	e.RangeComponents(func(component internal.Component) bool {
		if !component.(ComponentWhole).getStarted() {
			component.(ComponentWhole).setStarted(true)
			if cs, ok := component.(internal.ComponentStart); ok {
				cs.Start()
			}
		}
		return true
	})
}

func (e *Entity) CallUpdate() {
	if e.destroyed {
		return
	}

	e.RangeComponents(func(component internal.Component) bool {
		if component.(ComponentWhole).getStarted() {
			if cs, ok := component.(internal.ComponentUpdate); ok {
				cs.Update()
			}
		}
		return true
	})
}

func (e *Entity) CallLateUpdate() {
	if e.destroyed {
		return
	}

	e.RangeComponents(func(component internal.Component) bool {
		if component.(ComponentWhole).getStarted() {
			if cs, ok := component.(internal.ComponentLateUpdate); ok {
				cs.LateUpdate()
			}
		}
		return true
	})
}
