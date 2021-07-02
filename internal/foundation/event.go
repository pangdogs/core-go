package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal"
)

func BindEvent(hook, eventSrc interface{}, priority ...int) error {
	if hook == nil {
		return errors.New("nil hook")
	}

	if eventSrc == nil {
		return errors.New("nil eventSrc")
	}

	h := hook.(HookWhole)
	s := eventSrc.(EventSourceWhole)

	if err := h.attachEventSource(s); err != nil {
		return err
	}

	if err := s.addHook(h, priority...); err != nil {
		h.detachEventSource(s.GetEventSourceID())
		return err
	}

	return nil
}

func UnbindEvent(hook, eventSrc interface{}) {
	if hook == nil || eventSrc == nil {
		return
	}

	h := hook.(HookWhole)
	s := eventSrc.(EventSourceWhole)

	s.removeHook(h.GetHookID())
	h.detachEventSource(s.GetEventSourceID())
}

func UnbindAllEventSource(hook interface{}) {
	if hook == nil {
		return
	}

	hook.(HookWhole).rangeEventSources(func(eventSrc internal.EventSource) bool {
		UnbindEvent(hook, eventSrc)
		return true
	})
}

func UnbindAllHook(eventSrc interface{}) {
	if eventSrc == nil {
		return
	}

	eventSrc.(EventSourceWhole).rangeHooks(func(hook internal.Hook, priority int) bool {
		UnbindEvent(hook, eventSrc)
		return true
	})
}

func SendEvent(eventSrc interface{}, fun func(hook internal.Hook) bool) {
	if eventSrc == nil || fun == nil {
		return
	}

	eventSrc.(EventSourceWhole).rangeHooks(func(hook internal.Hook, priority int) bool {
		return fun(hook)
	})
}
