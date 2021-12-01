package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type Hook interface {
	InitHook(rt Runtime)
	GetHookID() uint64
	GetHookRuntime() Runtime
	SubscribeEvent(eventID int32, event misc.IFace) error
	UnsubscribeEvent(eventID int32)
	UnsubscribeAllEvent()
	GetEvent(eventID int32) misc.IFace
	addEventSource(eventSrc EventSource) (*misc.Element, error)
	removeEventSource(eventSrcEle *misc.Element)
	rangeEventSources(fun func(eventSrc EventSource) bool)
}

func IFace2Hook(f misc.IFace) Hook {
	return *(*Hook)(unsafe.Pointer(&f))
}

func Hook2IFace(h Hook) misc.IFace {
	return *(*misc.IFace)(unsafe.Pointer(&h))
}

var eventID = int32(-1)

func AllocEventID() int32 {
	eventID++

	if eventID >= eventsLimit {
		panic("eventID exceed limit")
	}

	return eventID
}

const eventsLimit = int32(64 * 3)

type HookFoundation struct {
	id             uint64
	runtime        Runtime
	eventSrcList   misc.List
	eventSrcGCList []*misc.Element
	eventList      misc.List
}

func (h *HookFoundation) GC() {
	for i := 0; i < len(h.eventSrcGCList); i++ {
		h.eventSrcList.Remove(h.eventSrcGCList[i])
	}
	h.eventSrcGCList = h.eventSrcGCList[:0]
}

func (h *HookFoundation) GCHandle() uintptr {
	return uintptr(unsafe.Pointer(h))
}

func (h *HookFoundation) InitHook(rt Runtime) {
	if rt == nil {
		panic("nil runtime")
	}

	h.id = rt.GetApp().makeUID()
	h.runtime = rt
	h.eventSrcList.Init(rt.GetCache())
	h.eventList.Init(rt.GetCache())
}

func (h *HookFoundation) GetHookID() uint64 {
	return h.id
}

func (h *HookFoundation) GetHookRuntime() Runtime {
	return h.runtime
}

func (h *HookFoundation) SubscribeEvent(eventID int32, event misc.IFace) error {
	if eventID < 0 || eventID >= eventsLimit {
		return errors.New("eventID invalid")
	}

	if event == misc.NilIFace {
		return errors.New("nil event")
	}

	if h.runtime == nil {
		return errors.New("nil runtime")
	}

	e := h.eventList.PushBack(nil)
	e.Mark[1] = uint64(eventID)
	h.runtime.subscribeEvent(h.id, eventID, event)

	return nil
}

func (h *HookFoundation) UnsubscribeEvent(eventID int32) {
	if eventID < 0 || eventID >= eventsLimit {
		return
	}

	if h.runtime == nil {
		return
	}

	h.runtime.unsubscribeEvent(h.id, eventID)
	h.eventList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Mark[1] == uint64(eventID) {
			h.eventList.Remove(e)
			return false
		}
		return true
	})
}

func (h *HookFoundation) UnsubscribeAllEvent() {
	h.eventList.UnsafeTraversal(func(e *misc.Element) bool {
		h.runtime.unsubscribeEvent(h.id, int32(e.Mark[1]))
		return true
	})
	h.eventList.Init(h.runtime.GetCache())
}

func (h *HookFoundation) GetEvent(eventID int32) misc.IFace {
	if eventID < 0 || eventID >= eventsLimit {
		panic("eventID invalid")
	}

	if h.runtime == nil {
		panic("nil runtime")
	}

	event, _ := h.runtime.getEvent(h.id, eventID)
	return event
}

func (h *HookFoundation) addEventSource(eventSrc EventSource) (*misc.Element, error) {
	if eventSrc == nil {
		return nil, errors.New("nil eventSrc")
	}

	if h.runtime == nil {
		return nil, errors.New("nil runtime")
	}

	eventSrcEle := h.eventSrcList.PushIFaceBack(EventSource2IFace(eventSrc))
	return eventSrcEle, nil
}

func (h *HookFoundation) removeEventSource(eventSrcEle *misc.Element) {
	if eventSrcEle == nil {
		return
	}

	eventSrcEle.SetMark(0, true)

	if h.runtime.GCEnabled() {
		h.eventSrcGCList = append(h.eventSrcGCList, eventSrcEle)
		h.runtime.PushGC(h)
	}
}

func (h *HookFoundation) rangeEventSources(fun func(eventSrc EventSource) bool) {
	if fun == nil {
		return
	}

	h.eventSrcList.UnsafeTraversal(func(ele *misc.Element) bool {
		if ele.Escape() || ele.GetMark(0) {
			return true
		}
		return fun(IFace2EventSource(ele.GetIFace(0)))
	})
}
