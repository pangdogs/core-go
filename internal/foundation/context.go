package foundation

import (
	"context"
	"github.com/pangdogs/core/internal"
	"sync"
)

func NewContext(parentCtx context.Context, reportError ...chan error) internal.Context {
	ctx := &Context{}
	ctx.InitContext(parentCtx, reportError...)
	return ctx
}

type Context struct {
	context.Context
	parentContext context.Context
	reportError   chan error
	cancel        context.CancelFunc
	wg            *sync.WaitGroup
	valueMap      sync.Map
}

func (ctx *Context) InitContext(parentCtx context.Context, reportError ...chan error) {
	if parentCtx == nil {
		ctx.parentContext = context.Background()
	} else {
		ctx.parentContext = parentCtx
	}

	if len(reportError) > 0 {
		ctx.reportError = reportError[0]
	}

	ctx.Context, ctx.cancel = context.WithCancel(ctx.parentContext)
	ctx.wg = &sync.WaitGroup{}
}

func (ctx *Context) GetParentContext() context.Context {
	return ctx.parentContext
}

func (ctx *Context) GetReportError() chan error {
	return ctx.reportError
}

func (ctx *Context) GetOrSetValue(key string, value interface{}) (actual interface{}, got bool) {
	return ctx.valueMap.LoadOrStore(key, value)
}

func (ctx *Context) SetValue(key string, value interface{}) {
	ctx.valueMap.Store(key, value)
}

func (ctx *Context) GetValue(key string) interface{} {
	value, _ := ctx.valueMap.Load(key)
	return value
}

func (ctx *Context) GetWaitGroup() *sync.WaitGroup {
	return ctx.wg
}

func (ctx *Context) GetCancelFunc() context.CancelFunc {
	return ctx.cancel
}
