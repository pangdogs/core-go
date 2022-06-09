package core

import (
	"github.com/pangdogs/core/container"
)

func BindEvent[T any](event IEvent, delegate T) Hook {
	return BindEventWithPriority(event, delegate, 0)
}

func BindEventWithPriority[T any](event IEvent, delegate T, priority int32) Hook {
	if event == nil {
		panic("nil event")
	}
	return event.newHook(delegate, IFace2Fast(delegate), priority)
}

type Hook struct {
	delegate          interface{}
	delegateFastIFace FastIFace
	priority          int32
	element           *container.Element[Hook]
}

func (hook *Hook) Bind(event IEvent) {
	hook.BindWithPriority(event, 0)
}

func (hook *Hook) BindWithPriority(event IEvent, priority int32) {
	if event == nil {
		panic("nil event")
	}

	if hook.element != nil && !hook.element.Escaped() {
		panic("repeated bind event invalid")
	}

	*hook = event.newHook(hook.delegate, hook.delegateFastIFace, priority)
}

func (hook *Hook) Unbind() {
	if hook.element != nil {
		hook.element.Escape()
		hook.element = nil
	}
}

func (hook *Hook) Delegate() FastIFace {
	return hook.delegateFastIFace
}
