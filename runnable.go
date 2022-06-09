package core

import "sync/atomic"

type Runnable interface {
	Run() chan struct{}
	Stop()
}

type RunnableBehavior struct {
	runningFlag int32
}

func (r *RunnableBehavior) markRunning() bool {
	return atomic.CompareAndSwapInt32(&r.runningFlag, 0, 1)
}

func (r *RunnableBehavior) markShutdown() bool {
	return atomic.CompareAndSwapInt32(&r.runningFlag, 1, 0)
}
