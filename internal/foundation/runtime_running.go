package foundation

import (
	"time"
)

func (rt *RuntimeFoundation) invokeEntityStart() {
	count := len(rt.entityStartList)
	if count <= 0 {
		return
	}

	for _, e := range rt.entityStartList {
		if e.Escape() || e.GetMark(0) {
			continue
		}
		CallOuter(rt.enableAutoRecover, rt.GetReportError(), IFace2Entity(e.GetIFace(0)).callStart)
	}

	rt.entityStartList = append(rt.entityStartList[:0], rt.entityStartList[count:]...)
}

func (rt *RuntimeFoundation) invokeEntityUpdate() {
	rt.frame.updateBegin()
	defer rt.frame.updateEnd()

	rt.RangeEntities(func(entity Entity) bool {
		entity.callUpdate()
		return true
	})

	rt.RangeEntities(func(entity Entity) bool {
		entity.callLateUpdate()
		return true
	})
}

func (rt *RuntimeFoundation) invokeSafeCall(callBundle *SafeCallBundle) (ret SafeRet) {
	ret.Err = CallOuter(rt.enableAutoRecover, rt.GetReportError(), func() {
		if callBundle.Stack != nil {
			ret = callBundle.SafeFun(callBundle.Stack)
		} else {
			ret = callBundle.UnsafeFun()
		}
	})
	return
}

func (rt *RuntimeFoundation) loopNoFrame() {
	CallOuter(rt.enableAutoRecover, rt.GetReportError(), func() {
		if rt.startFunc != nil {
			rt.startFunc(rt)
		}
	})

	for {
		for len(rt.entityStartList) > 0 {
			rt.invokeEntityStart()
		}

		select {
		case callBundle, ok := <-rt.safeCallList:
			if !ok {
				return
			}
			callBundle.Ret <- rt.invokeSafeCall(callBundle)

		case <-rt.Done():
			return
		}

		rt.RunGC()
	}
}

func (rt *RuntimeFoundation) loopWithFrame() {
	CallOuter(rt.enableAutoRecover, rt.GetReportError(), func() {
		rt.frame = rt.frameCreatorFunc(rt)
	})

	var ticker *time.Ticker

	if !rt.frame.IsBlink() {
		ticker = time.NewTicker(time.Duration(float64(time.Second) / float64(rt.frame.GetTargetFPS())))
		defer ticker.Stop()
	}

	CallOuter(rt.enableAutoRecover, rt.GetReportError(), func() {
		if rt.startFunc != nil {
			rt.startFunc(rt)
		}
	})

	rt.frame.cycleBegin()
	defer rt.frame.cycleEnd()

	for ; ; rt.frame.setCurFrames(rt.frame.GetCurFrames() + 1) {
		if !rt.loopWithFrameOnce(ticker) {
			return
		}
	}
}

func (rt *RuntimeFoundation) loopWithFrameOnce(ticker *time.Ticker) bool {
	if rt.frame.GetTotalFrames() > 0 {
		if rt.frame.GetCurFrames() >= rt.frame.GetTotalFrames() {
			return false
		}
	}

	rt.frame.frameBegin()
	defer rt.frame.frameEnd()

	onceUpdate := false

	if ticker != nil {
		for {
			for len(rt.entityStartList) > 0 {
				rt.invokeEntityStart()
			}

			select {
			case callBundle, ok := <-rt.safeCallList:
				if !ok {
					return false
				}
				callBundle.Ret <- rt.invokeSafeCall(callBundle)

			case <-ticker.C:
				if onceUpdate {
					return true
				}
				onceUpdate = true

				rt.invokeEntityUpdate()

			case <-rt.Done():
				return false
			}

			rt.RunGC()
		}
	} else {
		for {
			for len(rt.entityStartList) > 0 {
				rt.invokeEntityStart()
			}

			select {
			case callBundle, ok := <-rt.safeCallList:
				if !ok {
					return false
				}
				callBundle.Ret <- rt.invokeSafeCall(callBundle)

			case <-rt.Done():
				return false

			default:
				if onceUpdate {
					return true
				}
				onceUpdate = true

				rt.invokeEntityUpdate()
			}

			rt.RunGC()
		}
	}
}

func (rt *RuntimeFoundation) running(shutChan chan struct{}) {
	if parentCtx, ok := rt.GetParentContext().(Context); ok {
		parentCtx.GetWaitGroup().Add(1)
	}

	defer func() {
		close(rt.safeCallList)

	label:
		for {
			for len(rt.entityStartList) > 0 {
				rt.invokeEntityStart()
			}

			select {
			case callBundle, ok := <-rt.safeCallList:
				if !ok {
					break label
				}
				callBundle.Ret <- rt.invokeSafeCall(callBundle)

			default:
				break label
			}
		}

		if parentCtx, ok := rt.GetParentContext().(Context); ok {
			parentCtx.GetWaitGroup().Done()
		}

		rt.GetWaitGroup().Wait()
		rt.markShutdown()

		CallOuter(rt.enableAutoRecover, rt.GetReportError(), func() {
			if rt.stopFunc != nil {
				rt.stopFunc(rt)
			}
		})

		shutChan <- struct{}{}
	}()

	if rt.gcTimeInterval > 0 && rt.gcLastRunTime.IsZero() {
		rt.gcLastRunTime = time.Now()
	}

	rt.frame = nil

	if rt.frameCreatorFunc == nil {
		rt.loopNoFrame()
	} else {
		rt.loopWithFrame()
	}
}
