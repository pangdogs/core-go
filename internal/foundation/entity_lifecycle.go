package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type EntityLifecycleCaller interface {
	CallEntityInit()
	CallStart()
	CallUpdate()
	CallLateUpdate()
	CallEntityShut()
}

func IFace2EntityLifecycleCaller(f misc.IFace) EntityLifecycleCaller {
	return *(*EntityLifecycleCaller)(unsafe.Pointer(&f))
}

func EntityLifecycleCaller2IFace(elc EntityLifecycleCaller) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&elc))
}

const (
	EntityComponentsMark_Removed int = iota
	EntityComponentsMark_Inited
	EntityComponentsMark_Started
)

const (
	EntityComponentsIFace_Component int = iota
	EntityComponentsIFace_ComponentUpdate
	EntityComponentsIFace_ComponentLateUpdate
)

func (e *EntityFoundation) CallEntityInit() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if !e.GetMark(EntityComponentsMark_Inited) {
			e.SetMark(EntityComponentsMark_Inited, true)

			if cl, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentEntityInit); ok {
				cl.EntityInit()
			}
		}
		return true
	})
}

func (e *EntityFoundation) CallStart() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if !e.GetMark(EntityComponentsMark_Started) {
			e.SetMark(EntityComponentsMark_Started, true)

			if cl, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentStart); ok {
				cl.Start()
			}

			if cu, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentUpdate); ok {
				e.SetIFace(EntityComponentsIFace_ComponentUpdate, ComponentUpdate2IFace(cu))
			}

			if clu, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentLateUpdate); ok {
				e.SetIFace(EntityComponentsIFace_ComponentLateUpdate, ComponentLateUpdate2IFace(clu))
			}
		}
		return true
	})
}

func (e *EntityFoundation) CallUpdate() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Started) {
			return true
		}
		if cu := IFace2ComponentUpdate(e.GetIFace(EntityComponentsIFace_ComponentUpdate)); cu != nil {
			cu.Update()
		}
		return true
	})
}

func (e *EntityFoundation) CallLateUpdate() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Started) {
			return true
		}
		if clu := IFace2ComponentLateUpdate(e.GetIFace(EntityComponentsIFace_ComponentLateUpdate)); clu != nil {
			clu.LateUpdate()
		}
		return true
	})
}

func (e *EntityFoundation) CallEntityShut() {
	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Inited) {
			return true
		}
		if cl, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentEntityShut); ok {
			cl.EntityShut()
		}
		return true
	})
}
