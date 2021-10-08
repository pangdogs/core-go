package foundation

import (
	"errors"
)

func NewSafeCallBundle(stack SafeStack, safeFun func(stack SafeStack) SafeRet,
	unsafeFun func() SafeRet, retChan chan SafeRet) (*SafeCallBundle, error) {
	if stack != nil {
		if safeFun == nil {
			return nil, errors.New("nil safeFun")
		}
	} else {
		if unsafeFun == nil {
			return nil, errors.New("nil unsafeFun")
		}
	}

	if retChan == nil {
		return nil, errors.New("nil retChan")
	}

	return &SafeCallBundle{
		Stack:     stack,
		SafeFun:   safeFun,
		UnsafeFun: unsafeFun,
		Ret:       retChan,
	}, nil
}

type SafeCallBundle struct {
	Stack     SafeStack
	SafeFun   func(stack SafeStack) SafeRet
	UnsafeFun func() SafeRet
	Ret       chan SafeRet
}

type SafeRet struct {
	Err error
	Ret interface{}
}
