package foundation

import (
	"errors"
)

func BindEvent(hook, eventSrc interface{}, _priority ...int32) error {
	if hook == nil {
		return errors.New("nil hook")
	}

	if eventSrc == nil {
		return errors.New("nil eventSrc")
	}

	h := hook.(Hook)
	s := eventSrc.(EventSource)

	if err := h.attachEventSource(s); err != nil {
		return err
	}

	priority := int32(0)
	if len(_priority) > 0 {
		priority = _priority[0]
	}

	if err := s.addHook(h, priority); err != nil {
		h.detachEventSource(s.GetEventSourceID())
		return err
	}

	return nil
}

func UnbindEvent(hook, eventSrc interface{}) {
	if hook == nil || eventSrc == nil {
		return
	}

	h := hook.(Hook)
	s := eventSrc.(EventSource)

	s.removeHook(h.GetHookID())
	h.detachEventSource(s.GetEventSourceID())
}

func UnbindAllEventSource(hook interface{}) {
	if hook == nil {
		return
	}

	hook.(Hook).rangeEventSources(func(eventSrc EventSource) bool {
		UnbindEvent(hook, eventSrc)
		return true
	})
}

func UnbindAllHook(eventSrc interface{}) {
	if eventSrc == nil {
		return
	}

	eventSrc.(EventSource).rangeHooks(func(hook interface{}, priority int32) bool {
		UnbindEvent(hook, eventSrc)
		return true
	})
}

func SendEvent(eventSrc interface{}, fun func(hook interface{}) bool) {
	if eventSrc == nil || fun == nil {
		return
	}

	eventSrc.(EventSource).rangeHooks(func(hook interface{}, priority int32) bool {
		return fun(hook)
	})
}
