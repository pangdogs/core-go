package core

type SafeCall interface {
	SafeCall(segment func() SafeRet) <-chan SafeRet
	SafeCallNoRet(segment func())
	EventPushSafeCallSegment() IEvent
}

type SafeRet struct {
	Err error
	Ret interface{}
}
