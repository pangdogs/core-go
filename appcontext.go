package core

import (
	"context"
	"sync"
	"sync/atomic"
)

type AppContext interface {
	Context
	EntityMgr
	init(ctx context.Context, opts *AppContextOptions)
	getOptions() *AppContextOptions
	genUID() uint64
}

func AppContextGetOptions(appCtx AppContext) AppContextOptions {
	return *appCtx.getOptions()
}

func AppContextGetInheritor(appCtx AppContext) AppContext {
	return appCtx.getOptions().Inheritor
}

func NewAppContext(ctx context.Context, optFuncs ...NewAppContextOptionFunc) AppContext {
	opts := &AppContextOptions{}
	NewAppContextOption.Default()(opts)

	for _, optFunc := range optFuncs {
		optFunc(opts)
	}

	if opts.Inheritor != nil {
		opts.Inheritor.init(ctx, opts)
		return opts.Inheritor
	}

	app := &AppContextBehavior{}
	app.init(ctx, opts)

	return app.opts.Inheritor
}

type AppContextBehavior struct {
	ContextBehavior
	opts      AppContextOptions
	uidGen    uint64
	entityMap sync.Map
}

func (appCtx *AppContextBehavior) init(ctx context.Context, opts *AppContextOptions) {
	if ctx == nil {
		panic("nil ctx")
	}

	if opts == nil {
		panic("nil opts")
	}

	appCtx.opts = *opts

	if appCtx.opts.Inheritor == nil {
		appCtx.opts.Inheritor = appCtx
	}

	appCtx.ContextBehavior.init(ctx)
}

func (appCtx *AppContextBehavior) getOptions() *AppContextOptions {
	return &appCtx.opts
}

func (appCtx *AppContextBehavior) genUID() uint64 {
	return atomic.AddUint64(&appCtx.uidGen, 1)
}
