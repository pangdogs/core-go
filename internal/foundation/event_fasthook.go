package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/misc"
	"unsafe"
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
	eventsType [eventsLimit]unsafe.Pointer
	eventsData unsafe.Pointer
}

func (h *FastHookFoundation) SubscribeEvent(eventID int32, event misc.IFace) error {
	if eventID < 0 || eventID >= int32(len(h.eventsType)) {
		return errors.New("eventID invalid")
	}

	if event == misc.NilIFace {
		return errors.New("nil event")
	}

	h.eventsType[eventID] = event[0]
	h.eventsData = event[1]

	return nil
}

func (h *FastHookFoundation) GetEvent(eventID int32) misc.IFace {
	if eventID < 0 || eventID >= int32(len(h.eventsType)) {
		panic("eventID invalid")
	}

	return misc.IFace{h.eventsType[eventID], h.eventsData}
}
