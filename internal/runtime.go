package internal

type Runtime interface {
	Runnable
	Context
	GetRuntimeID() uint64
	GetApp() App
	GetFrame() Frame
	SafeCall(fun func() SafeRet) chan SafeRet
}
