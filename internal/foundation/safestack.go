package foundation

import (
	"fmt"
)

type SafeStack interface {
	SafeCall(callee Runtime, fun func(stack SafeStack) SafeRet) SafeRet
}

func NewSafeStack(caller Runtime) SafeStack {
	t := &_SafeStack{caller}
	return t
}

type _SafeStack []Runtime

func (stack *_SafeStack) Push(rt Runtime) *_SafeStack {
	*stack = append(*stack, rt)
	return stack
}

func (stack *_SafeStack) Copy() *_SafeStack {
	t := append(make(_SafeStack, 0, len(*stack)+1), *stack...)
	return &t
}

func (stack *_SafeStack) SafeCall(callee Runtime, fun func(stack SafeStack) SafeRet) (ret SafeRet) {
	defer func() {
		if info := recover(); info != nil {
			if err, ok := info.(error); ok {
				ret = SafeRet{Err: err}
			} else {
				ret = SafeRet{Err: fmt.Errorf("%v", info)}
			}
		}
	}()

	if callee == nil {
		panic("nil callee")
	}

	if fun == nil {
		panic("nil fun")
	}

	for i := 0; i <= len(*stack); i++ {
		if (*stack)[i].GetRuntimeID() == callee.GetRuntimeID() {
			ret = fun(stack.Copy())
			return
		}
	}

	retChan := make(chan SafeRet, 1)

	newStack := stack.Copy().Push(callee)

	callBundle, err := NewSafeCallBundle(newStack, fun, nil, retChan)
	if err != nil {
		panic(err)
	}

	callee.pushSafeCall(callBundle)

	ret = <-callBundle.Ret

	return
}
