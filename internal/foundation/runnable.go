package foundation

import "sync/atomic"

type Runnable interface {
	Run() chan struct{}
	Stop()
}

type _Runnable struct {
	shutChan    chan struct{}
	runningFlag int32
}

func (r *_Runnable) initRunnable() {
	r.shutChan = make(chan struct{}, 1)
}

func (r *_Runnable) markRunning() bool {
	return atomic.CompareAndSwapInt32(&r.runningFlag, 0, 1)
}

func (r *_Runnable) markShutdown() bool {
	return atomic.CompareAndSwapInt32(&r.runningFlag, 1, 0)
}
