package foundation

import (
	"github.com/pangdogs/core"
)

type FastHook interface {
	Hook
	SubscribeEvent(eventID int32, event core.IFace)
	GetEvent(eventID int32) core.IFace
}

type FastHookFoundation struct {
	core.HookFoundation
	events [64 * 3]core.IFace
}

func (h *FastHookFoundation) Conv2FastHook() FastHook {
	return h
}

func (h *FastHookFoundation) SubscribeEvent(eventID int32, event core.IFace) {
	if eventID < 0 || eventID >= int32(len(h.events)) {
		panic("eventID invalid")
	}

	if event == core.NilIFace {
		panic("nil event")
	}

	h.events[eventID] = event
}

func (h *FastHookFoundation) GetEvent(eventID int32) core.IFace {
	if eventID < 0 || eventID >= int32(len(h.events)) {
		panic("eventID invalid")
	}

	return h.events[eventID]
}
