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
	SubscribeEvent(eventID int32, subscriber misc.IFace) error
	UnsubscribeEvent(eventID int32)
	UnsubscribeAllEvent()
	GetEventSubscriber(eventID int32) misc.IFace
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

type HookFoundation struct {
	id                 uint64
	runtime            Runtime
	eventSrcList       misc.List
	eventSrcGCList     []*misc.Element
	eventBits          [misc.StoreMarkLimit - 1]uint64
	eventID            int32
	eventSubscriber    misc.IFace
	eventSubscriberMap map[int32]misc.IFace
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
}

func (h *HookFoundation) GetHookID() uint64 {
	return h.id
}

func (h *HookFoundation) GetHookRuntime() Runtime {
	return h.runtime
}

func (h *HookFoundation) SubscribeEvent(eventID int32, subscriber misc.IFace) error {
	if eventID < 0 || eventID >= eventsLimit {
		return errors.New("eventID invalid")
	}

	if subscriber == misc.NilIFace {
		return errors.New("nil event")
	}

	h.setEventMark(eventID, true)

	if h.eventSubscriber == misc.NilIFace {
		h.eventID = eventID
		h.eventSubscriber = subscriber

	} else if h.eventSubscriber != subscriber {
		if h.eventSubscriberMap == nil {
			h.eventSubscriberMap = map[int32]misc.IFace{}
		}

		h.eventSubscriberMap[h.eventID] = h.eventSubscriber
		h.eventSubscriberMap[eventID] = subscriber

		h.eventID = 0
		h.eventSubscriber = misc.NilIFace
	}

	return nil
}

func (h *HookFoundation) UnsubscribeEvent(eventID int32) {
	if eventID < 0 || eventID >= eventsLimit {
		return
	}

	h.setEventMark(eventID, false)

	if len(h.eventSubscriberMap) <= 0 {
		h.eventID = 0
		h.eventSubscriber = misc.NilIFace
		return
	}

	delete(h.eventSubscriberMap, eventID)
}

func (h *HookFoundation) UnsubscribeAllEvent() {
	for i := 0; i < len(h.eventBits); i++ {
		h.eventBits[i] = 0
	}
	h.eventID = 0
	h.eventSubscriber = misc.NilIFace
	h.eventSubscriberMap = nil
}

func (h *HookFoundation) GetEventSubscriber(eventID int32) misc.IFace {
	if eventID < 0 || eventID >= eventsLimit {
		panic("eventID invalid")
	}

	if h.runtime == nil {
		panic("nil runtime")
	}

	if !h.getEventMark(eventID) {
		return misc.NilIFace
	}

	if h.eventSubscriber != misc.NilIFace {
		return h.eventSubscriber
	}

	if h.eventSubscriberMap != nil {
		subscriber, ok := h.eventSubscriberMap[eventID]
		if ok {
			return subscriber
		}
	}

	panic("construct event failed")
}

func (h *HookFoundation) setEventMark(eventID int32, v bool) {
	if v {
		h.eventBits[eventID/64] |= 1 << eventID
	} else {
		h.eventBits[eventID/64] &= ^(1 << eventID)
	}
}

func (h *HookFoundation) getEventMark(eventID int32) bool {
	return (h.eventBits[eventID/64]>>eventID)&uint64(1) == 1
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
		return fun(IFace2EventSource(ele.GetIFace()))
	})
}
