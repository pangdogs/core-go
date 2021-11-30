package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
)

var eventID = int32(-1)

func AllocEventID() int32 {
	eventID++

	if eventID >= eventsLimit {
		panic("eventID exceed limit")
	}

	return eventID
}

const eventsLimit = int32(64 * 3)

type FastHookFoundation struct {
	HookFoundation
	events [eventsLimit]misc.IFace
}

func (h *FastHookFoundation) SubscribeEvent(eventID int32, event misc.IFace) error {
	if eventID < 0 || eventID >= int32(len(h.events)) {
		return errors.New("eventID invalid")
	}

	if event == misc.NilIFace {
		return errors.New("nil event")
	}

	h.events[eventID] = event

	return nil
}

func (h *FastHookFoundation) GetEvent(eventID int32) misc.IFace {
	if eventID < 0 || eventID >= int32(len(h.events)) {
		panic("eventID invalid")
	}

	return h.events[eventID]
}
