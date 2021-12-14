package foundation

import (
	"github.com/pangdogs/core/internal/misc"
)

const (
	EntityComponentsMark_Removed int = iota
	EntityComponentsMark_Inited
	EntityComponentsMark_Started
	EntityComponentsMark_Update
	EntityComponentsMark_LateUpdate
)

func (e *EntityFoundation) callEntityInit() {
	if e.destroyed {
		return
	}

	if e.lifecycleEntityInitFunc != nil {
		e.lifecycleEntityInitFunc()
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}

		if !e.GetMark(EntityComponentsMark_Inited) {
			e.SetMark(EntityComponentsMark_Inited, true)

			if cei, ok := IFace2Component(e.GetIFace()).(ComponentEntityInit); ok {
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

	if e.lifecycleStartFunc != nil {
		e.lifecycleStartFunc()
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}

		if !e.GetMark(EntityComponentsMark_Started) {
			e.SetMark(EntityComponentsMark_Started, true)

			if cs, ok := IFace2Component(e.GetIFace()).(ComponentStart); ok {
				cs.Start()
			}

			if _, ok := IFace2Component(e.GetIFace()).(ComponentUpdate); ok {
				e.SetMark(EntityComponentsMark_Update, true)
			}

			if _, ok := IFace2Component(e.GetIFace()).(ComponentLateUpdate); ok {
				e.SetMark(EntityComponentsMark_LateUpdate, true)
			}
		}

		return true
	})
}

func (e *EntityFoundation) callUpdate() {
	if e.destroyed {
		return
	}

	if e.lifecycleUpdateFunc != nil {
		if !e.lifecycleUpdateFunc() {
			return
		}
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Started) ||
			!e.GetMark(EntityComponentsMark_Update) {
			return true
		}

		if cu := IFace2ComponentUpdate(e.GetIFace()); cu != nil {
			cu.Update()
		}

		return true
	})
}

func (e *EntityFoundation) callLateUpdate() {
	if e.destroyed {
		return
	}

	if e.lifecycleLateUpdateFunc != nil {
		if !e.lifecycleLateUpdateFunc() {
			return
		}
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Started) ||
			!e.GetMark(EntityComponentsMark_LateUpdate) {
			return true
		}

		if clu := IFace2ComponentLateUpdate(e.GetIFace()); clu != nil {
			clu.LateUpdate()
		}

		return true
	})
}

func (e *EntityFoundation) callEntityShut() {
	if e.lifecycleEntityShutFunc != nil {
		e.lifecycleEntityShutFunc()
	}

	e.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) || !e.GetMark(EntityComponentsMark_Inited) {
			return true
		}

		if cs, ok := IFace2Component(e.GetIFace()).(ComponentEntityShut); ok {
			cs.EntityShut()
		}

		return true
	})
}

func (e *EntityFoundation) setLifecycleEntityInitFunc(fun func()) {
	e.lifecycleEntityInitFunc = fun
}

func (e *EntityFoundation) setLifecycleStartFunc(fun func()) {
	e.lifecycleStartFunc = fun
}

func (e *EntityFoundation) setLifecycleUpdateFunc(fun func() bool) {
	e.lifecycleUpdateFunc = fun
}

func (e *EntityFoundation) setLifecycleLateUpdateFunc(fun func() bool) {
	e.lifecycleLateUpdateFunc = fun
}

func (e *EntityFoundation) setLifecycleEntityShutFunc(fun func()) {
	e.lifecycleEntityShutFunc = fun
}

func (e *EntityFoundation) getLifecycleEntityInitFunc() func() {
	return e.lifecycleEntityInitFunc
}

func (e *EntityFoundation) getLifecycleStartFunc() func() {
	return e.lifecycleStartFunc
}

func (e *EntityFoundation) getLifecycleUpdateFunc() func() bool {
	return e.lifecycleUpdateFunc
}

func (e *EntityFoundation) getLifecycleLateUpdateFunc() func() bool {
	return e.lifecycleLateUpdateFunc
}

func (e *EntityFoundation) getLifecycleEntityShutFunc() func() {
	return e.lifecycleEntityShutFunc
}
