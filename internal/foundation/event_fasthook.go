package foundation

import (
	"github.com/pangdogs/core/internal/misc"
)

type FastHookFoundation struct {
	HookFoundation
	events [64 * 3]misc.IFace
}

func (h *FastHookFoundation) SubscribeEvent(eventID int32, event misc.IFace) {
	if eventID < 0 || eventID >= int32(len(h.events)) {
		panic("eventID invalid")
	}

	if event == misc.NilIFace {
		panic("nil event")
	}

	h.events[eventID] = event
}

func (h *FastHookFoundation) GetEvent(eventID int32) misc.IFace {
	if eventID < 0 || eventID >= int32(len(h.events)) {
		panic("eventID invalid")
	}

	return h.events[eventID]
}
