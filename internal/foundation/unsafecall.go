package foundation

import (
	"fmt"
)

func UnsafeCall(callee Runtime, fun func() SafeRet) (ret chan SafeRet) {
	ret = make(chan SafeRet, 1)

	defer func() {
		if info := recover(); info != nil {
			if err, ok := info.(error); ok {
				ret <- SafeRet{Err: err}
			} else {
				ret <- SafeRet{Err: fmt.Errorf("%v", info)}
			}
		}
	}()

	callBundle, err := NewSafeCallBundle(nil, nil, fun, ret)
	if err != nil {
		panic(err)
	}

	callee.pushSafeCall(callBundle)

	return
}
