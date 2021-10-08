package foundation

import (
	"errors"
)

type Hook interface {
	InitHook(rt Runtime)
	GetHookID() uint64
	getRuntime() Runtime
	attachEventSource(eventSrc EventSource) error
	detachEventSource(eventSrcID uint64)
	rangeEventSources(fun func(eventSrc EventSource) bool)
}

type HookFoundation struct {
	id           uint64
	runtime      Runtime
	eventSrcList []EventSource
}

func (h *HookFoundation) InitHook(rt Runtime) {
	h.id = rt.GetApp().makeUID()
	h.runtime = rt
}

func (h *HookFoundation) GetHookID() uint64 {
	return h.id
}

func (h *HookFoundation) getRuntime() Runtime {
	return h.runtime
}

func (h *HookFoundation) attachEventSource(eventSrc EventSource) error {
	if eventSrc == nil {
		return errors.New("nil eventSrc")
	}

	for i := 0; i < len(h.eventSrcList); i++ {
		if eventSrc.GetEventSourceID() == h.eventSrcList[i].GetEventSourceID() {
			return errors.New("event source id already exists")
		}
	}

	h.eventSrcList = append(h.eventSrcList, eventSrc)

	return nil
}

func (h *HookFoundation) detachEventSource(eventSrcID uint64) {
	for i := 0; i < len(h.eventSrcList); i++ {
		if eventSrcID == h.eventSrcList[i].GetEventSourceID() {
			h.eventSrcList = append(h.eventSrcList[:i], h.eventSrcList[i+1:]...)
			break
		}
	}
}

func (h *HookFoundation) rangeEventSources(fun func(eventSrc EventSource) bool) {
	if fun == nil {
		return
	}

	for _, es := range h.eventSrcList {
		if !fun(es) {
			break
		}
	}
}
