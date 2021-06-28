package internal

import (
	"context"
	"sync"
)

type Context interface {
	context.Context
	GetParentContext() context.Context
	GetOrSetValue(key string, value interface{}) (actual interface{}, got bool)
	SetValue(key string, value interface{})
	GetValue(key string) interface{}
	GetWaitGroup() *sync.WaitGroup
	GetCancelFunc() context.CancelFunc
}
