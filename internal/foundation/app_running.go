package foundation

func (app *AppFoundation) running(shutChan chan struct{}) {
	if parentCtx, ok := app.GetParentContext().(Context); ok {
		parentCtx.GetWaitGroup().Add(1)
	}

	defer func() {
		if parentCtx, ok := app.GetParentContext().(Context); ok {
			parentCtx.GetWaitGroup().Done()
		}

		app.GetWaitGroup().Wait()
		app.markShutdown()

		CallOuter(app.enableAutoRecover, app.GetReportError(), func() {
			if app.stopFunc != nil {
				app.stopFunc(app)
			}
		})

		shutChan <- struct{}{}
	}()

	CallOuter(app.enableAutoRecover, app.GetReportError(), func() {
		if app.startFunc != nil {
			app.startFunc(app)
		}
	})

	select {
	case <-app.Done():
		return
	}
}
