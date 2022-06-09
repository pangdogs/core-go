package core

import (
	container2 "github.com/pangdogs/core/container"
)

type IEvent interface {
	Emit(fun func(delegate FastIFace) bool)
	newHook(delegate interface{}, delegateFastIFace FastIFace, priority int32) Hook
}

type Event struct {
	subscribers container2.List[Hook]
	autoRecover bool
	reportError chan error
}

func (event *Event) Init(autoRecover bool, reportError chan error, hookCache *container2.Cache[Hook]) {
	event.autoRecover = autoRecover
	event.reportError = reportError
	event.subscribers.Init(hookCache)
}

func (event *Event) GC() {
	event.subscribers.GC()
}

func (event *Event) Emit(fun func(delegate FastIFace) bool) {
	if fun == nil {
		return
	}

	event.subscribers.Traversal(func(e *container2.Element[Hook]) bool {
		if e.Value.delegateFastIFace != NilFastIFace {
			ret, err := CallOuter(event.autoRecover, event.reportError, func() bool {
				return fun(e.Value.delegateFastIFace)
			})
			if err != nil {
				return true
			}
			return ret
		}
		return true
	})
}

func (event *Event) newHook(delegate interface{}, delegateFastIFace FastIFace, priority int32) Hook {
	hook := Hook{
		delegate:          delegate,
		delegateFastIFace: delegateFastIFace,
	}

	var mark *container2.Element[Hook]

	event.subscribers.ReverseTraversal(func(other *container2.Element[Hook]) bool {
		if hook.priority >= other.Value.priority {
			mark = other
			return false
		}
		return true
	})

	hook.priority = priority

	if mark != nil {
		hook.element = event.subscribers.InsertAfter(Hook{}, mark)
	} else {
		hook.element = event.subscribers.PushBack(Hook{})
	}

	return hook
}

func (event *Event) Clear() {
	event.subscribers.Traversal(func(e *container2.Element[Hook]) bool {
		e.Value.Unbind()
		return true
	})
}
