package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal"
)

func NewSafeCallBundle(stack internal.SafeStack, safeFun func(stack internal.SafeStack) internal.SafeRet,
	unsafeFun func() internal.SafeRet, retChan chan internal.SafeRet) (*SafeCallBundle, error) {
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
	Stack     internal.SafeStack
	SafeFun   func(stack internal.SafeStack) internal.SafeRet
	UnsafeFun func() internal.SafeRet
	Ret       chan internal.SafeRet
}
