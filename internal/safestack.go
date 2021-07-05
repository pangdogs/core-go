package internal

type SafeStack interface {
	SafeCall(callee Runtime, waitRet bool, fun func(stack SafeStack) SafeRet) SafeRet
}
