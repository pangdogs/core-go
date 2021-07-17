package foundation

import (
	"fmt"
	"github.com/pangdogs/core/internal"
)

func NewSafeStack(caller internal.Runtime) internal.SafeStack {
	t := &SafeStack{caller}
	return t
}

type SafeStack []internal.Runtime

func (stack *SafeStack) Push(rt internal.Runtime) *SafeStack {
	*stack = append(*stack, rt)
	return stack
}

func (stack *SafeStack) Copy() *SafeStack {
	t := append(make(SafeStack, 0, len(*stack)+1), *stack...)
	return &t
}

func (stack *SafeStack) SafeCall(callee internal.Runtime, fun func(stack internal.SafeStack) internal.SafeRet) (ret internal.SafeRet) {
	defer func() {
		if info := recover(); info != nil {
			if err, ok := info.(error); ok {
				ret = internal.SafeRet{Err: err}
			} else {
				ret = internal.SafeRet{Err: fmt.Errorf("%v", info)}
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

	retChan := make(chan internal.SafeRet, 1)

	newStack := stack.Copy().Push(callee)

	callBundle, err := NewSafeCallBundle(newStack, fun, nil, retChan)
	if err != nil {
		panic(err)
	}

	callee.(RuntimeWhole).PushSafeCall(callBundle)

	ret = <-callBundle.Ret

	return
}
