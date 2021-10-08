package foundation

import (
	"context"
	"sync"
)

type Context interface {
	context.Context
	GetParentContext() context.Context
	GetReportError() chan error
	GetOrSetValue(key string, value interface{}) (actual interface{}, got bool)
	SetValue(key string, value interface{})
	GetValue(key string) interface{}
	GetWaitGroup() *sync.WaitGroup
	GetCancelFunc() context.CancelFunc
}

func NewContext(parentCtx context.Context, reportError ...chan error) Context {
	ctx := &_Context{}
	ctx.initContext(parentCtx, reportError...)
	return ctx
}

type _Context struct {
	context.Context
	parentContext context.Context
	reportError   chan error
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	valueMap      sync.Map
}

func (ctx *_Context) initContext(parentCtx context.Context, reportError ...chan error) {
	if parentCtx == nil {
		ctx.parentContext = context.Background()
	} else {
		ctx.parentContext = parentCtx
	}

	if len(reportError) > 0 {
		ctx.reportError = reportError[0]
	}

	ctx.Context, ctx.cancel = context.WithCancel(ctx.parentContext)
}

func (ctx *_Context) GetParentContext() context.Context {
	return ctx.parentContext
}

func (ctx *_Context) GetReportError() chan error {
	return ctx.reportError
}

func (ctx *_Context) GetOrSetValue(key string, value interface{}) (actual interface{}, got bool) {
	return ctx.valueMap.LoadOrStore(key, value)
}

func (ctx *_Context) SetValue(key string, value interface{}) {
	ctx.valueMap.Store(key, value)
}

func (ctx *_Context) GetValue(key string) interface{} {
	value, _ := ctx.valueMap.Load(key)
	return value
}

func (ctx *_Context) GetWaitGroup() *sync.WaitGroup {
	return &ctx.wg
}

func (ctx *_Context) GetCancelFunc() context.CancelFunc {
	return ctx.cancel
}
