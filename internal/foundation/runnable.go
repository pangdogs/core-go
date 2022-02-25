package foundation

import "sync/atomic"

type Runnable interface {
	Run() chan struct{}
	Stop()
}

type RunnableFoundation struct {
	runningFlag int32
}

func (r *RunnableFoundation) markRunning() bool {
	return atomic.CompareAndSwapInt32(&r.runningFlag, 0, 1)
}

func (r *RunnableFoundation) markShutdown() bool {
	return atomic.CompareAndSwapInt32(&r.runningFlag, 1, 0)
}
