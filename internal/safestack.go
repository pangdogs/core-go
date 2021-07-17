package internal

type SafeStack interface {
	SafeCall(callee Runtime, fun func(stack SafeStack) SafeRet) SafeRet
}
