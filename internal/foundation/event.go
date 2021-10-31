package foundation

import (
	"errors"
	"unsafe"
)

func BindEvent(hook, eventSrc interface{}, _priority ...int32) error {
	if hook == nil {
		return errors.New("nil hook")
	}

	if eventSrc == nil {
		return errors.New("nil eventSrc")
	}

	h := hook.(Hook)
	es := eventSrc.(EventSource)

	rt := es.GetEventSourceRuntime()
	if rt == nil {
		return errors.New("nil runtime")
	}

	if rt.eventIsBound(h.GetHookID(), es.GetEventSourceID()) {
		return errors.New("already bound")
	}

	priority := int32(0)
	if len(_priority) > 0 {
		priority = _priority[0]
	}

	hookEle, err := es.addHook(h, priority)
	if err != nil {
		return err
	}

	eventSrcEle, err := h.addEventSource(es)
	if err != nil {
		es.removeHook(hookEle)
		return err
	}

	if err := rt.bindEvent(h.GetHookID(), es.GetEventSourceID(), hookEle, eventSrcEle); err != nil {
		es.removeHook(hookEle)
		h.removeEventSource(eventSrcEle)
		return err
	}

	return nil
}

func UnbindEvent(hook, eventSrc interface{}) {
	if hook == nil || eventSrc == nil {
		return
	}

	unbindEvent(hook.(Hook), eventSrc.(EventSource))
}

func unbindEvent(hook Hook, eventSrc EventSource) {
	rt := eventSrc.GetEventSourceRuntime()
	if rt == nil {
		return
	}

	hookEle, eventSrcEle, ok := rt.unbindEvent(hook.GetHookID(), eventSrc.GetEventSourceID())
	if !ok {
		return
	}

	eventSrc.removeHook(hookEle)
	hook.removeEventSource(eventSrcEle)
}

func UnbindAllEventSource(hook interface{}) {
	if hook == nil {
		return
	}

	h := hook.(Hook)

	h.rangeEventSources(func(eventSrc interface{}) bool {
		unbindEvent(h, eventSrc.(EventSource))
		return true
	})
}

func UnbindAllHook(eventSrc interface{}) {
	if eventSrc == nil {
		return
	}

	es := eventSrc.(EventSource)

	es.rangeHooks(func(hook interface{}, priority int32) bool {
		unbindEvent(hook.(Hook), es)
		return true
	})
}

type EventRet int32

const (
	EventRet_Continue EventRet = 1 << iota
	EventRet_Break
	EventRet_Unsubscribe
)

func SendEvent(eventSrc interface{}, fun func(hook interface{}) EventRet) {
	if eventSrc == nil || fun == nil {
		return
	}

	eventSrc.(EventSource).sendEvent(fun, **(**uintptr)(unsafe.Pointer(&fun)))
}
