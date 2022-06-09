package core

import "time"

func (runtime *RuntimeBehavior) Run() chan struct{} {
	if !runtime.markRunning() {
		panic("runtime already running")
	}

	shutChan := make(chan struct{}, 1)

	go runtime.running(shutChan)

	return shutChan
}

func (runtime *RuntimeBehavior) Stop() {
	runtime.ctx.GetCancelFunc()()
}

func (runtime *RuntimeBehavior) OnPushSafeCallSegment(segment func()) {
	timer := time.NewTimer(runtime.opts.ProcessQueueTimeout)
	defer timer.Stop()

	select {
	case runtime.processQueue <- segment:
		return
	case <-timer.C:
		panic("process queue push segment timeout")
	}
}

func (runtime *RuntimeBehavior) running(shutChan chan struct{}) {
	if parentCtx, ok := runtime.ctx.GetParentCtx().(Context); ok {
		parentCtx.GetWaitGroup().Add(1)
	}

	hooks := runtime.loopStarted()

	defer func() {
		runtime.loopStopped(hooks)

		if parentCtx, ok := runtime.ctx.GetParentCtx().(Context); ok {
			parentCtx.GetWaitGroup().Done()
		}

		runtime.ctx.GetWaitGroup().Wait()

		runtime.markShutdown()
		shutChan <- struct{}{}
	}()

	frame := runtime.opts.Frame

	if frame == nil {
		defer runtime.loopNoFrameEnd()
		runtime.loopNoFrame()

	} else if frame.Blink() {
		defer runtime.loopWithBlinkFrameEnd()
		runtime.loopWithBlinkFrame()

	} else {
		defer runtime.loopWithFrameEnd()
		runtime.loopWithFrame()
	}
}

func (runtime *RuntimeBehavior) loopStarted() (hooks [5]Hook) {
	runtimeCtx := runtime.ctx
	frame := runtime.opts.Frame

	runtime.processQueue = make(chan func(), runtime.opts.ProcessQueueCapacity)

	if frame != nil {
		frame.runningBegin()
	}

	hooks[0] = BindEvent[EventEntityMgrAddEntity[RuntimeContext]](runtimeCtx.EventEntityMgrAddEntity(), runtime)
	hooks[1] = BindEvent[EventEntityMgrRemoveEntity[RuntimeContext]](runtimeCtx.EventEntityMgrRemoveEntity(), runtime)
	hooks[2] = BindEvent[EventEntityMgrEntityAddComponents[RuntimeContext]](runtimeCtx.EventEntityMgrEntityAddComponents(), runtime)
	hooks[3] = BindEvent[EventEntityMgrEntityRemoveComponent[RuntimeContext]](runtimeCtx.EventEntityMgrEntityRemoveComponent(), runtime)
	hooks[4] = BindEvent[EventPushSafeCallSegment](runtimeCtx.EventPushSafeCallSegment(), runtime)

	runtimeCtx.RangeEntities(func(entity Entity) bool {
		CallOuterNoRet(runtime.opts.EnableAutoRecover, runtimeCtx.GetReportError(), func() {
			runtime.OnEntityMgrAddEntity(runtimeCtx, entity)
		})
		return true
	})

	CallOuterNoRet(runtime.opts.EnableAutoRecover, runtimeCtx.GetReportError(), func() {
		if runtimeCtx.getOptions().StartedCallback != nil {
			runtimeCtx.getOptions().StartedCallback(runtime.opts.Inheritor)
		}
	})

	return
}

func (runtime *RuntimeBehavior) loopStopped(hooks [5]Hook) {
	runtimeCtx := runtime.ctx
	frame := runtime.opts.Frame

	CallOuterNoRet(runtime.opts.EnableAutoRecover, runtimeCtx.GetReportError(), func() {
		if runtimeCtx.getOptions().StoppedCallback != nil {
			runtimeCtx.getOptions().StoppedCallback(runtime.opts.Inheritor)
		}
	})

	runtimeCtx.RangeEntities(func(entity Entity) bool {
		CallOuterNoRet(runtime.opts.EnableAutoRecover, runtimeCtx.GetReportError(), func() {
			runtime.OnEntityMgrRemoveEntity(runtimeCtx, entity)
		})
		return true
	})

	for _, hook := range hooks {
		hook.Unbind()
	}

	if frame != nil {
		frame.runningEnd()
	}
}

func (runtime *RuntimeBehavior) loopNoFrame() {
	gcTicker := time.NewTicker(runtime.opts.GCInterval)
	defer gcTicker.Stop()

	for {
		select {
		case process, ok := <-runtime.processQueue:
			if ok {
				CallOuterNoRet(runtime.opts.EnableAutoRecover, runtime.ctx.GetReportError(), process)
			}

		case <-gcTicker.C:
			runtime.GC()

		case <-runtime.ctx.Done():
			return
		}
	}
}

func (runtime *RuntimeBehavior) loopNoFrameEnd() {
	close(runtime.processQueue)

	for {
		select {
		case process, ok := <-runtime.processQueue:
			if ok {
				CallOuterNoRet(runtime.opts.EnableAutoRecover, runtime.ctx.GetReportError(), process)
			}

		default:
			return
		}
	}
}

func (runtime *RuntimeBehavior) loopWithFrame() {
	frame := runtime.opts.Frame

	ticker := time.NewTicker(time.Duration(float64(time.Second) / float64(frame.GetTargetFPS())))
	defer ticker.Stop()

	go func() {
		totalFrames := frame.GetTotalFrames()

		for i := uint64(0); ; i++ {
			if totalFrames > 0 && i >= totalFrames {
				return
			}

			select {
			case <-ticker.C:
				CallOuterNoRet(runtime.opts.EnableAutoRecover, runtime.ctx.GetReportError(), func() {
					timer := time.NewTimer(runtime.opts.ProcessQueueTimeout)
					defer timer.Stop()

					select {
					case runtime.processQueue <- runtime.frameUpdate:
						return
					case <-timer.C:
						panic("process queue push frame update timeout")
					}
				})
			}
		}
	}()

	gcTicker := time.NewTicker(runtime.opts.GCInterval)
	defer gcTicker.Stop()

	for {
		select {
		case process, ok := <-runtime.processQueue:
			if ok {
				CallOuterNoRet(runtime.opts.EnableAutoRecover, runtime.ctx.GetReportError(), process)
			}

		case <-gcTicker.C:
			runtime.GC()

		case <-runtime.ctx.Done():
			return
		}
	}
}

func (runtime *RuntimeBehavior) loopWithFrameEnd() {
	frame := runtime.opts.Frame

	close(runtime.processQueue)

	for {
		select {
		case process, ok := <-runtime.processQueue:
			if ok {
				CallOuterNoRet(runtime.opts.EnableAutoRecover, runtime.ctx.GetReportError(), process)
			}

		default:
			break
		}
	}

	if frame.GetCurFrames() > 0 {
		frame.frameEnd()
	}
}

func (runtime *RuntimeBehavior) frameUpdate() {
	frame := runtime.opts.Frame

	if frame.GetCurFrames() > 0 {
		frame.frameEnd()
	}
	frame.frameBegin()

	frame.updateBegin()
	defer frame.updateEnd()

	emitEventUpdate(&runtime.eventUpdate)
	emitEventLateUpdate(&runtime.eventLateUpdate)
}

func (runtime *RuntimeBehavior) loopWithBlinkFrame() {
	frame := runtime.opts.Frame
	totalFrames := frame.GetTotalFrames()

	gcFrames := uint64(runtime.opts.GCInterval.Seconds() * float64(frame.GetTargetFPS()))

	for frame.setCurFrames(0); ; {
		curFrames := frame.GetCurFrames()

		if totalFrames > 0 && curFrames >= totalFrames {
			return
		}

		if !runtime.blinkFrameUpdate() {
			return
		}

		if curFrames%gcFrames == 0 {
			runtime.GC()
		}

		frame.setCurFrames(curFrames + 1)
	}
}

func (runtime *RuntimeBehavior) loopWithBlinkFrameEnd() {
	close(runtime.processQueue)

	for {
		select {
		case process, ok := <-runtime.processQueue:
			if ok {
				CallOuterNoRet(runtime.opts.EnableAutoRecover, runtime.ctx.GetReportError(), process)
			}

		default:
			break
		}
	}
}

func (runtime *RuntimeBehavior) blinkFrameUpdate() bool {
	frame := runtime.opts.Frame

	frame.frameBegin()
	defer frame.frameEnd()

	for {
		select {
		case process, ok := <-runtime.processQueue:
			if ok {
				CallOuterNoRet(runtime.opts.EnableAutoRecover, runtime.ctx.GetReportError(), process)
			}

		case <-runtime.ctx.Done():
			return false

		default:
			func() {
				frame.updateBegin()
				defer frame.updateEnd()

				emitEventUpdate(&runtime.eventUpdate)
				emitEventLateUpdate(&runtime.eventLateUpdate)
			}()
		}
	}

	return true
}
