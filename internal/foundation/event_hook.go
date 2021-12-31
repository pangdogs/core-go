package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

type Hook interface {
	initHook(rt Runtime)
	GetHookID() uint64
	GetHookRuntime() Runtime
	SubscribeEvent(eventID int32, subscriber misc.IFace) error
	UnsubscribeEvent(eventID int32)
	UnsubscribeAllEvent()
	GetEventSubscriber(eventID int32) misc.IFace
	setEventCalled(eventID int32, v bool)
	getEventCalled(eventID int32) bool
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

func InitHook(hook Hook, rt Runtime) {
	hook.initHook(rt)
}

type HookFoundation struct {
	id                  uint64
	runtime             Runtime
	eventSrcList        misc.List
	eventSrcGCList      []*misc.Element
	eventSubscriberMap  map[int32]misc.IFace
	eventSubscribedMark [eventsLimit / 64]uint64
	eventCalledMark     [eventsLimit / 64]uint64
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

func (h *HookFoundation) initHook(rt Runtime) {
	if rt == nil {
		panic("nil runtime")
	}

	if h.runtime != nil {
		panic("init repeated")
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

	if h.eventSubscriberMap == nil {
		h.eventSubscriberMap = map[int32]misc.IFace{}
	}

	h.eventSubscriberMap[eventID] = subscriber
	h.setEventSubscribed(eventID, true)

	return nil
}

func (h *HookFoundation) UnsubscribeEvent(eventID int32) {
	if eventID < 0 || eventID >= eventsLimit {
		return
	}

	if h.eventSubscriberMap == nil {
		return
	}

	delete(h.eventSubscriberMap, eventID)
	h.setEventSubscribed(eventID, false)
}

func (h *HookFoundation) UnsubscribeAllEvent() {
	h.eventSubscriberMap = nil
	for i := range h.eventSubscribedMark {
		h.eventSubscribedMark[i] = 0
	}
}

func (h *HookFoundation) GetEventSubscriber(eventID int32) misc.IFace {
	if eventID < 0 || eventID >= eventsLimit {
		panic("eventID invalid")
	}

	if h.eventSubscriberMap == nil {
		return misc.NilIFace
	}

	if !h.getEventSubscribed(eventID) {
		return misc.NilIFace
	}

	subscriber, ok := h.eventSubscriberMap[eventID]
	if !ok {
		return misc.NilIFace
	}

	return subscriber
}

func (h *HookFoundation) setEventSubscribed(eventID int32, v bool) {
	if v {
		h.eventSubscribedMark[eventID/64] |= 1 << (eventID % 64)
	} else {
		h.eventSubscribedMark[eventID/64] &= ^(1 << (eventID % 64))
	}
}

func (h *HookFoundation) getEventSubscribed(eventID int32) bool {
	return (h.eventSubscribedMark[eventID/64]>>(eventID%64))&1 != 0
}

func (h *HookFoundation) setEventCalled(eventID int32, v bool) {
	if v {
		h.eventCalledMark[eventID/64] |= 1 << (eventID % 64)
	} else {
		h.eventCalledMark[eventID/64] &= ^(1 << (eventID % 64))
	}
}

func (h *HookFoundation) getEventCalled(eventID int32) bool {
	return (h.eventCalledMark[eventID/64]>>(eventID%64))&1 != 0
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

	h.eventSrcList.UnsafeTraversal(func(e *misc.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}
		return fun(IFace2EventSource(e.GetIFace(0)))
	})
}
