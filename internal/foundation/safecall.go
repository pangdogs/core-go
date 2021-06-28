package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal"
)

func NewSafeCallBundle(fun func() internal.SafeRet) (*SafeCallBundle, error) {
	if fun == nil {
		return nil, errors.New("nil fun")
	}

	return &SafeCallBundle{
		Fun: fun,
		Ret: make(chan internal.SafeRet, 1),
	}, nil
}

type SafeCallBundle struct {
	Fun func() internal.SafeRet
	Ret chan internal.SafeRet
}
