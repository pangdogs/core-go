package foundation

import (
	"errors"
)

func BindEvent(hook, eventSrc interface{}, priority ...int) (ret error) {
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

	defer func() {
		if ret != nil {
			h.detachEventSource(s.GetEventSourceID())
		}
	}()

	if err := s.addHook(h, priority...); err != nil {
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

	eventSrc.(EventSource).rangeHooks(func(hook Hook, priority int) bool {
		UnbindEvent(hook, eventSrc)
		return true
	})
}

func SendEvent(eventSrc interface{}, fun func(hook Hook) bool) {
	if eventSrc == nil || fun == nil {
		return
	}

	eventSrc.(EventSource).rangeHooks(func(hook Hook, priority int) bool {
		return fun(hook)
	})
}
