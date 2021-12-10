package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type EventSource interface {
	InitEventSource(rt Runtime)
	GetEventSourceID() uint64
	GetEventSourceRuntime() Runtime
	addHook(hook Hook, priority int32) (*misc.Element, error)
	removeHook(hookEle *misc.Element)
	rangeHooks(fun func(hook Hook, priority int32) bool)
	sendEvent(fun func(hook Hook) EventRet, eventHandle uintptr)
}

func IFace2EventSource(f misc.IFace) EventSource {
	return *(*EventSource)(unsafe.Pointer(&f))
}

func EventSource2IFace(es EventSource) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&es))
}

type EventSourceFoundation struct {
	id         uint64
	runtime    Runtime
	hookList   misc.List
	hookGCList []*misc.Element
}

func (es *EventSourceFoundation) GC() {
	for i := 0; i < len(es.hookGCList); i++ {
		es.hookList.Remove(es.hookGCList[i])
	}
	es.hookGCList = es.hookGCList[:0]
}

func (es *EventSourceFoundation) GCHandle() uintptr {
	return uintptr(unsafe.Pointer(es))
}

func (es *EventSourceFoundation) InitEventSource(rt Runtime) {
	if rt == nil {
		panic("nil runtime")
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

	for ele := es.hookList.Front(); ele != nil; ele = ele.Next() {
		if priority < int32(ele.Mark[0]>>32) {
			hookEle := es.hookList.InsertIFaceBefore(Hook2IFace(hook), ele)
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

	es.hookList.UnsafeTraversal(func(ele *misc.Element) bool {
		if ele.Escape() || ele.GetMark(0) {
			return true
		}
		return fun(IFace2Hook(ele.GetIFace()), int32(ele.Mark[0]>>32))
	})
}

func (es *EventSourceFoundation) sendEvent(fun func(hook Hook) EventRet, eventHandle uintptr) {
	if fun == nil {
		return
	}

	bit := es.runtime.eventHandleToBit(eventHandle)

	es.hookList.UnsafeTraversal(func(ele *misc.Element) bool {
		if ele.Escape() || ele.GetMark(0) || ele.GetMark(bit) {
			return true
		}

		ret := fun(IFace2Hook(ele.GetIFace()))

		if ret&EventRet_Unsubscribe != 0 {
			ele.SetMark(bit, true)
		}

		return ret&EventRet_Break == 0
	})
}
