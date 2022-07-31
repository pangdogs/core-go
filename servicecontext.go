package core

import (
	"context"
	"sync"
	"sync/atomic"
)

type ServiceContext interface {
	Context
	_RunnableMark
	EntityMgr
	init(ctx context.Context, opts *ServiceContextOptions)
	getOptions() *ServiceContextOptions
	genUID() uint64
}

func ServiceContextGetOptions(servCtx ServiceContext) ServiceContextOptions {
	return *servCtx.getOptions()
}

func ServiceContextGetInheritor(servCtx ServiceContext) ServiceContext {
	return servCtx.getOptions().Inheritor
}

func NewServiceContext(ctx context.Context, optFuncs ...NewServiceContextOptionFunc) ServiceContext {
	opts := &ServiceContextOptions{}
	NewServiceContextOption.Default()(opts)

	for i := range optFuncs {
		optFuncs[i](opts)
	}

	if opts.Inheritor != nil {
		opts.Inheritor.init(ctx, opts)
		return opts.Inheritor
	}

	serv := &ServiceContextBehavior{}
	serv.init(ctx, opts)

	return serv.opts.Inheritor
}

type ServiceContextBehavior struct {
	_ContextBehavior
	_RunnableMarkBehavior
	opts      ServiceContextOptions
	uidGen    uint64
	entityMap sync.Map
}

func (servCtx *ServiceContextBehavior) init(ctx context.Context, opts *ServiceContextOptions) {
	if ctx == nil {
		panic("nil ctx")
	}

	if opts == nil {
		panic("nil opts")
	}

	servCtx.opts = *opts

	if servCtx.opts.Inheritor == nil {
		servCtx.opts.Inheritor = servCtx
	}

	servCtx._ContextBehavior.init(ctx, servCtx.opts.ReportError)
}

func (servCtx *ServiceContextBehavior) getOptions() *ServiceContextOptions {
	return &servCtx.opts
}

func (servCtx *ServiceContextBehavior) genUID() uint64 {
	return atomic.AddUint64(&servCtx.uidGen, 1)
}
