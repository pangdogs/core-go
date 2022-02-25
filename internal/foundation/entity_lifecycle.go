package foundation

import (
	"github.com/pangdogs/core/internal/misc"
)

const (
	EntityComponentsMark_Removed int = iota
	EntityComponentsMark_Inited
	EntityComponentsMark_Awaked
	EntityComponentsMark_EntityInited
	EntityComponentsMark_Started
	EntityComponentsMark_Update
	EntityComponentsMark_LateUpdate
	EntityComponentsMark_EntityShut
	EntityComponentsMark_Halted
	EntityComponentsMark_Shut
)

func (ent *EntityFoundation) callEntityInit() {
	if ent.destroying {
		return
	}

	if ent.lifecycleEntityInitFunc != nil {
		ent.lifecycleEntityInitFunc()
	}

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}

		if !e.GetMark(EntityComponentsMark_EntityInited) {
			e.SetMark(EntityComponentsMark_EntityInited, true)

			if cei, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentEntityInit); ok {
				cei.EntityInit()
			}
		}

		return true
	})
}

func (ent *EntityFoundation) callStart() {
	if ent.destroying {
		return
	}

	if ent.lifecycleStartFunc != nil {
		ent.lifecycleStartFunc()
	}

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}

		if !e.GetMark(EntityComponentsMark_Started) {
			e.SetMark(EntityComponentsMark_Started, true)

			if cs, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentStart); ok {
				cs.Start()
			}

			if cu, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentUpdate); ok {
				e.SetIFace(EntityComponentsIFace_ComponentUpdate, ComponentUpdate2IFace(cu))
				e.SetMark(EntityComponentsMark_Update, true)
				ent.componentsLifecycleUpdateCount++
			}

			if clu, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentLateUpdate); ok {
				e.SetIFace(EntityComponentsIFace_ComponentLateUpdate, ComponentLateUpdate2IFace(clu))
				e.SetMark(EntityComponentsMark_LateUpdate, true)
				ent.componentsLifecycleLateUpdateCount++
			}
		}

		return true
	})
}

func (ent *EntityFoundation) callUpdate() {
	if ent.destroying {
		return
	}

	if ent.lifecycleUpdateFunc != nil {
		if !ent.lifecycleUpdateFunc() {
			return
		}
	}

	if ent.componentsLifecycleUpdateCount <= 0 {
		return
	}

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) ||
			!e.GetMark(EntityComponentsMark_Started) || !e.GetMark(EntityComponentsMark_Update) {
			return true
		}

		if cu := IFace2ComponentUpdate(e.GetIFace(EntityComponentsIFace_ComponentUpdate)); cu != nil {
			cu.Update()
		}

		return true
	})
}

func (ent *EntityFoundation) callLateUpdate() {
	if ent.destroying {
		return
	}

	if ent.lifecycleLateUpdateFunc != nil {
		if !ent.lifecycleLateUpdateFunc() {
			return
		}
	}

	if ent.componentsLifecycleLateUpdateCount <= 0 {
		return
	}

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) ||
			!e.GetMark(EntityComponentsMark_Started) || !e.GetMark(EntityComponentsMark_LateUpdate) {
			return true
		}

		if clu := IFace2ComponentLateUpdate(e.GetIFace(EntityComponentsIFace_ComponentLateUpdate)); clu != nil {
			clu.LateUpdate()
		}

		return true
	})
}

func (ent *EntityFoundation) callEntityShut() {
	if ent.lifecycleEntityShutFunc != nil {
		ent.lifecycleEntityShutFunc()
	}

	ent.componentList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) ||
			!e.GetMark(EntityComponentsMark_EntityInited) {
			return true
		}

		if !e.GetMark(EntityComponentsMark_EntityShut) {
			e.SetMark(EntityComponentsMark_EntityShut, true)

			if cs, ok := IFace2Component(e.GetIFace(EntityComponentsIFace_Component)).(ComponentEntityShut); ok {
				cs.EntityShut()
			}
		}

		return true
	})
}

func (ent *EntityFoundation) setLifecycleEntityInitFunc(fun func()) {
	ent.lifecycleEntityInitFunc = fun
}

func (ent *EntityFoundation) setLifecycleStartFunc(fun func()) {
	ent.lifecycleStartFunc = fun
}

func (ent *EntityFoundation) setLifecycleUpdateFunc(fun func() bool) {
	ent.lifecycleUpdateFunc = fun
}

func (ent *EntityFoundation) setLifecycleLateUpdateFunc(fun func() bool) {
	ent.lifecycleLateUpdateFunc = fun
}

func (ent *EntityFoundation) setLifecycleEntityShutFunc(fun func()) {
	ent.lifecycleEntityShutFunc = fun
}

func (ent *EntityFoundation) getLifecycleEntityInitFunc() func() {
	return ent.lifecycleEntityInitFunc
}

func (ent *EntityFoundation) getLifecycleStartFunc() func() {
	return ent.lifecycleStartFunc
}

func (ent *EntityFoundation) getLifecycleUpdateFunc() func() bool {
	return ent.lifecycleUpdateFunc
}

func (ent *EntityFoundation) getLifecycleLateUpdateFunc() func() bool {
	return ent.lifecycleLateUpdateFunc
}

func (ent *EntityFoundation) getLifecycleEntityShutFunc() func() {
	return ent.lifecycleEntityShutFunc
}
