package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal"
)

type HookWhole interface {
	internal.Hook
	attachEventSource(eventSrc internal.EventSource) error
	detachEventSource(eventSrcID uint64)
	rangeEventSources(fun func(eventSrc internal.EventSource) bool)
}

type Hook struct {
	id           uint64
	eventSrcList []internal.EventSource
}

func (h *Hook) InitHook(rt internal.Runtime) {
	h.id = rt.GetApp().(AppWhole).MakeUID()
}

func (h *Hook) BindEvent(eventSrc EventSource) {

}

func (h *Hook) GetHookID() uint64 {
	return h.id
}

func (h *Hook) attachEventSource(eventSrc internal.EventSource) error {
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

func (h *Hook) detachEventSource(eventSrcID uint64) {
	for i := 0; i < len(h.eventSrcList); i++ {
		if eventSrcID == h.eventSrcList[i].GetEventSourceID() {
			h.eventSrcList = append(h.eventSrcList[:i], h.eventSrcList[i+1:]...)
			break
		}
	}
}

func (h *Hook) rangeEventSources(fun func(eventSrc internal.EventSource) bool) {
	if fun == nil {
		return
	}

	for _, es := range h.eventSrcList {
		if !fun(es) {
			break
		}
	}
}
