package foundation

import (
	"sync"
	"sync/atomic"
)

type App interface {
	Runnable
	Context
	initApp(ctx Context, opts *AppOptions)
	getInheritor() App
	GetEntity(entID uint64) Entity
	RangeEntities(func(entity Entity) bool)
	makeUID() uint64
	addEntity(entity Entity)
	removeEntity(entID uint64)
}

func AppGetInheritor(app App) App {
	return app.getInheritor()
}

func NewApp(ctx Context, optFuncs ...NewAppOptionFunc) App {
	opts := &AppOptions{}
	NewAppOption.Default()(opts)

	for _, optFun := range optFuncs {
		optFun(opts)
	}

	if opts.inheritor != nil {
		opts.inheritor.initApp(ctx, opts)
		return opts.inheritor
	}

	app := &AppFoundation{}
	app.initApp(ctx, opts)

	return app.inheritor
}

type AppFoundation struct {
	RunnableFoundation
	Context
	AppOptions
	uidMaker  uint64
	entityMap sync.Map
}

func (app *AppFoundation) initApp(ctx Context, opts *AppOptions) {
	if ctx == nil {
		panic("nil ctx")
	}

	if opts == nil {
		panic("nil opts")
	}

	app.AppOptions = *opts

	if app.inheritor == nil {
		app.inheritor = app
	}

	app.initRunnable()
	app.Context = ctx

	CallOuter(app.autoRecover, app.GetReportError(), func() {
		if app.initFunc != nil {
			app.initFunc(app)
		}
	})
}

func (app *AppFoundation) Run() chan struct{} {
	if !app.markRunning() {
		panic("app already running")
	}

	go func() {
		if parentCtx, ok := app.GetParentContext().(Context); ok {
			parentCtx.GetWaitGroup().Add(1)
		}

		defer func() {
			if parentCtx, ok := app.GetParentContext().(Context); ok {
				parentCtx.GetWaitGroup().Done()
			}

			app.GetWaitGroup().Wait()
			app.markShutdown()
			app.shutChan <- struct{}{}

			CallOuter(app.autoRecover, app.GetReportError(), func() {
				if app.stopFunc != nil {
					app.stopFunc(app)
				}
			})
		}()

		CallOuter(app.autoRecover, app.GetReportError(), func() {
			if app.startFunc != nil {
				app.startFunc(app)
			}
		})

		select {
		case <-app.Done():
			return
		}
	}()

	return app.shutChan
}

func (app *AppFoundation) Stop() {
	app.GetCancelFunc()()
}

func (app *AppFoundation) getInheritor() App {
	return app.inheritor
}

func (app *AppFoundation) GetEntity(entID uint64) Entity {
	if !app.enableGetEntity {
		return nil
	}

	entity, ok := app.entityMap.Load(entID)
	if !ok {
		return nil
	}

	return entity.(Entity)
}

func (app *AppFoundation) RangeEntities(fun func(entity Entity) bool) {
	if fun == nil || !app.enableGetEntity {
		return
	}

	app.entityMap.Range(func(key, value interface{}) bool {
		return fun(value.(Entity))
	})
}

func (app *AppFoundation) makeUID() uint64 {
	return atomic.AddUint64(&app.uidMaker, 1)
}

func (app *AppFoundation) addEntity(entity Entity) {
	if !app.enableGetEntity {
		return
	}

	if entity == nil {
		panic("nil entity")
	}

	if _, loaded := app.entityMap.LoadOrStore(entity.GetEntityID(), entity); loaded {
		panic("entity id already exists")
	}
}

func (app *AppFoundation) removeEntity(entID uint64) {
	if !app.enableGetEntity {
		return
	}

	app.entityMap.Delete(entID)
}
