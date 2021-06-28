package internal

type Hook interface {
	GetHookID() uint64
	InitHook(rt Runtime)
}

type EventSource interface {
	GetEventSourceID() uint64
	InitEventSource(rt Runtime)
}
