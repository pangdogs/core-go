package core

import (
	"errors"
	"github.com/pangdogs/core/container"
)

func (entity *EntityBehavior) GetComponent(name string) Component {
	if e, ok := entity.getComponentElement(name); ok {
		return Fast2IFace[Component](e.Value.FastIFace)
	}

	return nil
}

func (entity *EntityBehavior) GetComponentByID(id uint64) Component {
	if e, ok := entity.getComponentElementByID(id); ok {
		return Fast2IFace[Component](e.Value.FastIFace)
	}

	return nil
}

func (entity *EntityBehavior) GetComponents(name string) []Component {
	if e, ok := entity.getComponentElement(name); ok {
		var components []Component

		entity.componentList.TraversalAt(func(other *container.Element[Face]) bool {
			comp := Fast2IFace[Component](other.Value.FastIFace)
			if comp.GetName() == name {
				components = append(components, comp)
				return true
			}
			return false
		}, e)

		return components
	}

	return nil
}

func (entity *EntityBehavior) RangeComponents(fun func(component Component) bool) {
	if fun == nil {
		return
	}

	entity.componentList.Traversal(func(e *container.Element[Face]) bool {
		return fun(Fast2IFace[Component](e.Value.FastIFace))
	})
}

func (entity *EntityBehavior) AddComponents(name string, components []Component) error {
	for _, comp := range components {
		if err := entity.addSingleComponent(name, comp); err != nil {
			return err
		}
	}

	emitEventCompMgrAddComponents(&entity.eventCompMgrAddComponents, entity.opts.Inheritor, components)
	return nil
}

func (entity *EntityBehavior) AddComponent(name string, component Component) error {
	if err := entity.addSingleComponent(name, component); err != nil {
		return err
	}

	emitEventCompMgrAddComponents(&entity.eventCompMgrAddComponents, entity.opts.Inheritor, []Component{component})
	return nil
}

func (entity *EntityBehavior) RemoveComponent(name string) {
	e, ok := entity.getComponentElement(name)
	if !ok {
		return
	}

	if entity.opts.EnableFastGetComponent {
		delete(entity.componentMap, name)
	}

	entity.componentList.TraversalAt(func(other *container.Element[Face]) bool {
		comp := Fast2IFace[Component](other.Value.FastIFace)
		if comp.GetName() == name {
			other.Escape()
			emitEventCompMgrRemoveComponent(&entity.eventCompMgrRemoveComponent, entity.opts.Inheritor, comp)
			return true
		}
		return false
	}, e)
}

func (entity *EntityBehavior) RemoveComponentByID(id uint64) {
	e, ok := entity.getComponentElementByID(id)
	if !ok {
		return
	}

	if entity.opts.EnableFastGetComponentByID {
		delete(entity.componentByIDMap, id)
	}

	e.Escape()
	emitEventCompMgrRemoveComponent(&entity.eventCompMgrRemoveComponent, entity.opts.Inheritor, Fast2IFace[Component](e.Value.FastIFace))
}

func (entity *EntityBehavior) EventCompMgrAddComponents() IEvent {
	return &entity.eventCompMgrAddComponents
}

func (entity *EntityBehavior) EventCompMgrRemoveComponent() IEvent {
	return &entity.eventCompMgrRemoveComponent
}

func (entity *EntityBehavior) addSingleComponent(name string, component Component) error {
	if component == nil {
		return errors.New("nil component")
	}

	if component.GetEntity() != nil {
		return errors.New("component already added in entity")
	}

	component.init(name, entity.opts.Inheritor, component, entity.opts.HookCache)

	face := Face{
		IFace:     component,
		FastIFace: IFace2Fast(component),
	}

	if e, ok := entity.getComponentElement(name); ok {
		entity.componentList.TraversalAt(func(other *container.Element[Face]) bool {
			if Fast2IFace[Component](other.Value.FastIFace).GetName() == name {
				e = other
				return true
			}
			return false
		}, e)

		e = entity.componentList.InsertAfter(face, e)
		e.GC = component

	} else {
		e = entity.componentList.PushBack(face)
		e.GC = component

		if entity.opts.EnableFastGetComponent {
			entity.componentMap[name] = e
		}

		if entity.opts.EnableFastGetComponentByID {
			entity.componentByIDMap[component.GetID()] = e
		}
	}

	return nil
}

func (entity *EntityBehavior) getComponentElement(name string) (*container.Element[Face], bool) {
	if entity.opts.EnableFastGetComponent {
		e, ok := entity.componentMap[name]
		return e, ok
	}

	var e *container.Element[Face]

	entity.componentList.Traversal(func(other *container.Element[Face]) bool {
		if Fast2IFace[Component](other.Value.FastIFace).GetName() == name {
			e = other
			return false
		}
		return true
	})

	return e, e != nil
}

func (entity *EntityBehavior) getComponentElementByID(id uint64) (*container.Element[Face], bool) {
	if entity.opts.EnableFastGetComponentByID {
		e, ok := entity.componentByIDMap[id]
		return e, ok
	}

	var e *container.Element[Face]

	entity.componentList.Traversal(func(other *container.Element[Face]) bool {
		if Fast2IFace[Component](other.Value.FastIFace).GetID() == id {
			e = other
			return false
		}
		return true
	})

	return e, e != nil
}
