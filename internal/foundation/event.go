package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
)

const eventsLimit = int32(64 * (misc.StoreMarkLimit - 1))

var eventID = int32(-1)

func AllocEventID() int32 {
	eventID++

	if eventID >= eventsLimit {
		panic("eventID exceed limit")
	}

	return eventID
}

func BindEvent(hook Hook, eventSrc EventSource, priority ...int32) error {
	if hook == nil {
		return errors.New("nil hook")
	}

	if hook.GetHookRuntime() == nil {
		return errors.New("nil hook runtime")
	}

	if eventSrc == nil {
		return errors.New("nil eventSrc")
	}

	rt := eventSrc.GetEventSourceRuntime()
	if rt == nil {
		return errors.New("nil eventSrc runtime")
	}

	if rt.eventIsBound(hook.GetHookID(), eventSrc.GetEventSourceID()) {
		return errors.New("hook and eventSrc already bound")
	}

	_priority := int32(0)
	if len(priority) > 0 {
		_priority = priority[0]
	}

	hookEle, err := eventSrc.addHook(hook, _priority)
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

func UnbindEvent(hook Hook, eventSrc EventSource) {
	if hook == nil || eventSrc == nil {
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

func SendEvent(eventSrc EventSource, fun func(hook Hook) EventRet) {
	if eventSrc == nil || fun == nil {
		return
	}

	eventSrc.sendEvent(fun, **(**uintptr)(unsafe.Pointer(&fun)))
}
