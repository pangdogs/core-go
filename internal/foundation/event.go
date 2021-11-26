package foundation

import (
	"errors"
	"unsafe"
)

func BindEvent(hook Hook, _eventSrc interface{}, _priority ...int32) error {
	if hook == nil {
		return errors.New("nil hook")
	}

	eventSrc, ok := _eventSrc.(EventSource)
	if !ok {
		return errors.New("eventSrc invalid")
	}

	rt := eventSrc.GetEventSourceRuntime()
	if rt == nil {
		return errors.New("nil runtime")
	}

	if rt.eventIsBound(hook.GetHookID(), eventSrc.GetEventSourceID()) {
		return errors.New("already bound")
	}

	priority := int32(0)
	if len(_priority) > 0 {
		priority = _priority[0]
	}

	hookEle, err := eventSrc.addHook(hook, priority)
	if err != nil {
		return err
	}

	eventSrcEle, err := hook.addEventSource(eventSrc)
	if err != nil {
		eventSrc.removeHook(hookEle)
		return err
	}

	if err := rt.bindEvent(hook.GetHookID(), eventSrc.GetEventSourceID(), hookEle, eventSrcEle); err != nil {
		eventSrc.removeHook(hookEle)
		hook.removeEventSource(eventSrcEle)
		return err
	}

	return nil
}

func UnbindEvent(hook Hook, _eventSrc interface{}) {
	if hook == nil {
		return
	}

	eventSrc, ok := _eventSrc.(EventSource)
	if !ok {
		return
	}

	unbindEvent(hook, eventSrc)
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

func UnbindAllEventSource(hook Hook) {
	if hook == nil {
		return
	}

	hook.rangeEventSources(func(eventSrc EventSource) bool {
		unbindEvent(hook, eventSrc)
		return true
	})
}

func UnbindAllHook(eventSrc EventSource) {
	if eventSrc == nil {
		return
	}

	eventSrc.rangeHooks(func(hook Hook, priority int32) bool {
		unbindEvent(hook, eventSrc)
		return true
	})
}

type EventRet int32

const (
	EventRet_Continue EventRet = 1 << iota
	EventRet_Break
	EventRet_Unsubscribe
)

func SendEvent(_eventSrc interface{}, fun func(hook Hook) EventRet) {
	if fun == nil {
		return
	}

	eventSrc, ok := _eventSrc.(EventSource)
	if !ok {
		return
	}

	eventSrc.sendEvent(fun, **(**uintptr)(unsafe.Pointer(&fun)))
}
