package core

func (app *AppBehavior) Run() chan struct{} {
	if !app.ctx.markRunning() {
		panic("app already running")
	}

	shutChan := make(chan struct{}, 1)

	go app.running(shutChan)

	return shutChan
}

func (app *AppBehavior) Stop() {
	app.ctx.GetCancelFunc()()
}

func (app *AppBehavior) running(shutChan chan struct{}) {
	if parentCtx, ok := app.ctx.GetParentCtx().(Context); ok {
		parentCtx.GetWaitGroup().Add(1)
	}

	defer func() {
		CallOuterNoRet(app.opts.EnableAutoRecover, app.ctx.GetReportError(), func() {
			if app.ctx.getOptions().StoppingCallback != nil {
				app.ctx.getOptions().StoppingCallback(app)
			}
		})

		if parentCtx, ok := app.ctx.GetParentCtx().(Context); ok {
			parentCtx.GetWaitGroup().Done()
		}

		app.ctx.GetWaitGroup().Wait()

		CallOuterNoRet(app.opts.EnableAutoRecover, app.ctx.GetReportError(), func() {
			if app.ctx.getOptions().StoppedCallback != nil {
				app.ctx.getOptions().StoppedCallback(app)
			}
		})

		app.ctx.markShutdown()
		shutChan <- struct{}{}
	}()

	CallOuterNoRet(app.opts.EnableAutoRecover, app.ctx.GetReportError(), func() {
		if app.ctx.getOptions().StartedCallback != nil {
			app.ctx.getOptions().StartedCallback(app)
		}
	})

	select {
	case <-app.ctx.Done():
		return
	}
}
