package foundation

import "github.com/pangdogs/core/internal/list"

type EntityLifecycleCaller interface {
	CallEntityInit()
	CallStart()
	CallUpdate()
	CallLateUpdate()
	CallEntityShut()
}

const (
	EntityComponentsMark_Removed int = iota
	EntityComponentsMark_Inited
	EntityComponentsMark_Started
	EntityComponentsMark_NoUpdate
	EntityComponentsMark_NoLateUpdate
)

func (e *EntityFoundation) CallEntityInit() {
	if e.destroyed {
		return
	}

	e.componentList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(EntityComponentsMark_Removed) {
			return true
		}
		if !e.GetMark(EntityComponentsMark_Inited) {
			e.SetMark(EntityComponentsMark_Inited, true)

			if cl, ok := e.Value.(ComponentEntityInit); ok {
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
		if e.GetMark(EntityComponentsMark_Inited) {
			if cl, ok := e.Value.(ComponentEntityShut); ok {
				cl.EntityShut()
			}
		}
		return true
	})
}
