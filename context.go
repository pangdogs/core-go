package core

import (
	"context"
	"sync"
)

type Context interface {
	context.Context
	GetParentCtx() context.Context
	GetReportError() chan error
	GetOrSetValue(key string, value interface{}) (actual interface{}, got bool)
	SetValue(key string, value interface{})
	GetValue(key string) (interface{}, bool)
	GetWaitGroup() *sync.WaitGroup
	GetCancelFunc() context.CancelFunc
}

func NewContext(parentCtx context.Context, reportError ...chan error) Context {
	ctx := &ContextBehavior{}
	ctx.init(parentCtx, reportError...)
	return ctx
}

type ContextBehavior struct {
	context.Context
	parentCtx   context.Context
	reportError chan error
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	valueMap    sync.Map
}

func (ctx *ContextBehavior) init(parentCtx context.Context, reportError ...chan error) {
	if parentCtx == nil {
		ctx.parentCtx = context.Background()
	} else {
		ctx.parentCtx = parentCtx
	}

	if len(reportError) > 0 {
		ctx.reportError = reportError[0]
	}

	ctx.Context, ctx.cancel = context.WithCancel(ctx.parentCtx)
}

func (ctx *ContextBehavior) GetParentCtx() context.Context {
	return ctx.parentCtx
}

func (ctx *ContextBehavior) GetReportError() chan error {
	return ctx.reportError
}

func (ctx *ContextBehavior) GetOrSetValue(key string, value interface{}) (actual interface{}, got bool) {
	return ctx.valueMap.LoadOrStore(key, value)
}

func (ctx *ContextBehavior) SetValue(key string, value interface{}) {
	ctx.valueMap.Store(key, value)
}

func (ctx *ContextBehavior) GetValue(key string) (interface{}, bool) {
	return ctx.valueMap.Load(key)
}

func (ctx *ContextBehavior) GetWaitGroup() *sync.WaitGroup {
	return &ctx.wg
}

func (ctx *ContextBehavior) GetCancelFunc() context.CancelFunc {
	return ctx.cancel
}
