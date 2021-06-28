package foundation

import "sync/atomic"

type Runnable struct {
	shutChan    chan struct{}
	runningFlag int32
}

func (r *Runnable) InitRunnable() {
	r.shutChan = make(chan struct{}, 1)
}

func (r *Runnable) MarkRunning() bool {
	return atomic.CompareAndSwapInt32(&r.runningFlag, 0, 1)
}

func (r *Runnable) MarkShutdown() bool {
	return atomic.CompareAndSwapInt32(&r.runningFlag, 1, 0)
}
