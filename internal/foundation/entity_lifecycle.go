package foundation

import (
	"github.com/pangdogs/core/internal/misc"
)

const (
	EntityComponentsMark_Removed int = iota
	EntityComponentsMark_Inited
	EntityComponentsMark_Started
)

func (e *EntityFoundation) callEntityInit() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if !e.GetMark(EntityComponentsMark_Inited) {
			e.SetMark(EntityComponentsMark_Inited, true)

			if cei := IFace2Component(e.GetIFace()).getLifecycleComponentEntityInit(); cei != nil {
				cei.EntityInit()
			}
		}
		return true
	})
}

func (e *EntityFoundation) callStart() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if !e.GetMark(EntityComponentsMark_Started) {
			e.SetMark(EntityComponentsMark_Started, true)

			if cs := IFace2Component(e.GetIFace()).getLifecycleComponentStart(); cs != nil {
				cs.Start()
			}
		}
		return true
	})
}

func (e *EntityFoundation) callUpdate() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Started) {
			return true
		}
		if cu := IFace2Component(e.GetIFace()).getLifecycleComponentUpdate(); cu != nil {
			cu.Update()
		}
		return true
	})
}

func (e *EntityFoundation) callLateUpdate() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Started) {
			return true
		}
		if clu := IFace2Component(e.GetIFace()).getLifecycleComponentLateUpdate(); clu != nil {
			clu.LateUpdate()
		}
		return true
	})
}

func (e *EntityFoundation) callEntityShut() {
	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Inited) {
			return true
		}
		if cs := IFace2Component(e.GetIFace()).getLifecycleComponentEntityShut(); cs != nil {
			cs.EntityShut()
		}
		return true
	})
}

func (e *EntityFoundation) CallEntityInit() {
	e.callEntityInit()
}

func (e *EntityFoundation) CallStart() {
	e.callStart()
}

func (e *EntityFoundation) CallUpdate() {
	e.callUpdate()
}

func (e *EntityFoundation) CallLateUpdate() {
	e.callLateUpdate()
}

func (e *EntityFoundation) CallEntityShut() {
	e.callEntityShut()
}
