package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type EventSource interface {
	initEventSource(rt Runtime)
	GetEventSourceID() uint64
	GetEventSourceRuntime() Runtime
	addHook(hook Hook, priority int32) (*misc.Element, error)
	removeHook(hookEle *misc.Element)
	rangeHooks(fun func(hook Hook, priority int32) bool)
	sendEvent(eventID int32, fun func(subscriber misc.IFace) EventRet)
}

func IFace2EventSource(f misc.IFace) EventSource {
	return *(*EventSource)(unsafe.Pointer(&f))
}

func EventSource2IFace(es EventSource) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&es))
}

func InitEventSource(eventSrc EventSource, rt Runtime) {
	eventSrc.initEventSource(rt)
}

type EventSourceFoundation struct {
	id         uint64
	runtime    Runtime
	hookList   misc.List
	hookGCList []*misc.Element
}

func (es *EventSourceFoundation) GC() {
	for i := range es.hookGCList {
		es.hookList.Remove(es.hookGCList[i])
	}
	es.hookGCList = es.hookGCList[:0]
}

func (es *EventSourceFoundation) GCHandle() uintptr {
	return uintptr(unsafe.Pointer(es))
}

func (es *EventSourceFoundation) initEventSource(rt Runtime) {
	if rt == nil {
		panic("nil runtime")
	}

	if es.runtime != nil {
		panic("init repeated")
	}

	es.id = rt.GetApp().makeUID()
	es.runtime = rt
	es.hookList.Init(rt.GetCache())
}

func (es *EventSourceFoundation) GetEventSourceID() uint64 {
	return es.id
}

func (es *EventSourceFoundation) GetEventSourceRuntime() Runtime {
	return es.runtime
}

func (es *EventSourceFoundation) addHook(hook Hook, priority int32) (*misc.Element, error) {
	if hook == nil {
		return nil, errors.New("nil hook")
	}

	if es.runtime == nil {
		return nil, errors.New("nil runtime")
	}

	for e := es.hookList.Front(); e != nil; e = e.Next() {
		if priority < int32(e.Mark[0]>>32) {
			hookEle := es.hookList.InsertIFaceBefore(Hook2IFace(hook), e)
			hookEle.Mark[0] |= uint64(priority) << 32
			return hookEle, nil
		}
	}

	hookEle := es.hookList.PushIFaceBack(Hook2IFace(hook))
	hookEle.Mark[0] |= uint64(priority) << 32
	return hookEle, nil
}

func (es *EventSourceFoundation) removeHook(hookEle *misc.Element) {
	if hookEle == nil {
		return
	}

	hookEle.SetMark(0, true)

	if es.runtime.GCEnabled() {
		es.hookGCList = append(es.hookGCList, hookEle)
		es.runtime.PushGC(es)
	}
}

func (es *EventSourceFoundation) rangeHooks(fun func(hook Hook, priority int32) bool) {
	if fun == nil {
		return
	}

	es.hookList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}
		return fun(IFace2Hook(e.GetIFace(0)), int32(e.Mark[0]>>32))
	})
}

func (es *EventSourceFoundation) sendEvent(eventID int32, fun func(subscriber misc.IFace) EventRet) {
	if fun == nil || es.runtime == nil {
		return
	}

	enableEventRecursion := es.runtime.eventRecursionEnabled()
	discardRecursiveEvent := es.runtime.recursiveEventDiscarded()

	es.hookList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}

		hook := IFace2Hook(e.GetIFace(0))

		subscriber := hook.GetEventSubscriber(eventID)
		if subscriber == misc.NilIFace {
			return true
		}

		called := hook.getEventCalled(eventID)
		if called {
			if enableEventRecursion {
				if discardRecursiveEvent {
					return true
				}
			} else {
				panic("event recursion not enabled")
			}
		} else {
			hook.setEventCalled(eventID, true)
			defer hook.setEventCalled(eventID, false)
		}

		return fun(subscriber) == EventRet_Continue
	})
}
