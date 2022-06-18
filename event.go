package core

import "github.com/pangdogs/core/container"

type IEvent interface {
	Emit(fun func(delegate FastIFace) bool)
	newHook(delegate interface{}, delegateFastIFace FastIFace, priority int32) Hook
}

type EventRecursion int32

const (
	EventRecursion_Allow EventRecursion = iota
	EventRecursion_Disallow
	EventRecursion_Discard
)

type Event struct {
	subscribers    container.List[Hook]
	autoRecover    bool
	reportError    chan error
	eventRecursion EventRecursion
	emitted        int
	opened         bool
}

func (event *Event) Init(autoRecover bool, reportError chan error, eventRecursion EventRecursion, hookCache *container.Cache[Hook], gcCollector container.GCCollector) {
	event.autoRecover = autoRecover
	event.reportError = reportError
	event.eventRecursion = eventRecursion
	event.subscribers.Init(hookCache, gcCollector)
	event.opened = true
}

func (event *Event) GC() {
	event.subscribers.GC()
}

func (event *Event) MarkGC() {
	event.subscribers.MarkGC()
}

func (event *Event) NeedGC() bool {
	return event.subscribers.NeedGC()
}

func (event *Event) Emit(fun func(delegate FastIFace) bool) {
	if fun == nil {
		return
	}

	if event.emitted > 0 {
		switch event.eventRecursion {
		case EventRecursion_Allow:
			break
		case EventRecursion_Disallow:
			panic("recursive event disallowed")
		case EventRecursion_Discard:
			return
		}
	}

	event.emitted++
	defer func() {
		event.emitted--
	}()

	event.subscribers.Traversal(func(e *container.Element[Hook]) bool {
		if !event.opened {
			return false
		}
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
	if !event.opened {
		panic("event closed")
	}

	hook := Hook{
		delegate:          delegate,
		delegateFastIFace: delegateFastIFace,
	}

	var mark *container.Element[Hook]

	event.subscribers.ReverseTraversal(func(other *container.Element[Hook]) bool {
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

	hook.element.Value = hook

	return hook
}

func (event *Event) Open() {
	event.opened = true
}

func (event *Event) Close() {
	event.Clear()
	event.opened = false
}

func (event *Event) Clear() {
	event.subscribers.Traversal(func(e *container.Element[Hook]) bool {
		e.Value.Unbind()
		return true
	})
}
