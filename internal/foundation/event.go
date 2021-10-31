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

	h := hook.(Hook)
	es := eventSrc.(EventSource)

	rt := es.GetEventSourceRuntime()
	if rt == nil {
		return
	}

	hookEle, eventSrcEle, ok := rt.unbindEvent(h.GetHookID(), es.GetEventSourceID())
	if !ok {
		return
	}

	es.removeHook(hookEle)
	h.removeEventSource(eventSrcEle)
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
