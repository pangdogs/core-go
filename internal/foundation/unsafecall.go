package foundation

import (
	"fmt"
	"github.com/pangdogs/core/internal"
)

func UnsafeCall(callee internal.Runtime, fun func() internal.SafeRet) (ret chan internal.SafeRet) {
	ret = make(chan internal.SafeRet, 1)

	defer func() {
		if info := recover(); info != nil {
			if err, ok := info.(error); ok {
				ret <- internal.SafeRet{Err: err}
			} else {
				ret <- internal.SafeRet{Err: fmt.Errorf("%v", info)}
			}
		}
	}()

	callBundle, err := NewSafeCallBundle(nil, nil, fun, ret)
	if err != nil {
		panic(err)
	}

	callee.(RuntimeWhole).PushSafeCall(callBundle)

	return
}
