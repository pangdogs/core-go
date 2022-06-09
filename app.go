package core

type App interface {
	Runnable
	init(ctx AppContext, opts *AppOptions)
	getOptions() *AppOptions
	GetContext() AppContext
}

func AppGetOptions(app App) AppOptions {
	return *app.getOptions()
}

func AppGetInheritor(app App) App {
	return app.getOptions().Inheritor
}

func NewApp(appCtx AppContext, optFuncs ...NewAppOptionFunc) App {
	opts := &AppOptions{}
	NewAppOption.Default()(opts)

	for _, optFunc := range optFuncs {
		optFunc(opts)
	}

	if opts.Inheritor != nil {
		opts.Inheritor.init(appCtx, opts)
		return opts.Inheritor
	}

	app := &AppBehavior{}
	app.init(appCtx, opts)

	return app.opts.Inheritor
}

type AppBehavior struct {
	RunnableBehavior
	opts AppOptions
	ctx  AppContext
}

func (app *AppBehavior) init(appCtx AppContext, opts *AppOptions) {
	if appCtx == nil {
		panic("nil appCtx")
	}

	if opts == nil {
		panic("nil opts")
	}

	app.opts = *opts

	if app.opts.Inheritor == nil {
		app.opts.Inheritor = app
	}

	app.ctx = appCtx
}

func (app *AppBehavior) getOptions() *AppOptions {
	return &app.opts
}

func (app *AppBehavior) GetContext() AppContext {
	return app.ctx
}
