package foundation

import (
	"github.com/pangdogs/core/internal/misc"
	"time"
)

var NewRuntimeOption = &NewRuntimeOptions{}

type RuntimeOptions struct {
	inheritor Runtime
	initFunc,
	startFunc,
	stopFunc func(rt Runtime)
	enableAutoRun         bool
	enableAutoRecover     bool
	safeCallCacheSize     int
	frameCreatorFunc      func(rt Runtime) Frame
	enableGC              bool
	gcTimeInterval        time.Duration
	gcItemNum             int
	cache                 *misc.Cache
	enableEventRecursion  bool
	discardRecursiveEvent bool
	callEventDepth        int
}

type NewRuntimeOptionFunc func(o *RuntimeOptions)

type NewRuntimeOptions struct{}

func (*NewRuntimeOptions) Default() NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.inheritor = nil
		o.initFunc = nil
		o.startFunc = nil
		o.stopFunc = nil
		o.enableAutoRun = false
		o.enableAutoRecover = false
		o.safeCallCacheSize = 128
		o.frameCreatorFunc = nil
		o.enableGC = true
		o.gcTimeInterval = 0
		o.gcItemNum = 512
		o.cache = nil
		o.enableEventRecursion = false
		o.discardRecursiveEvent = true
		o.callEventDepth = 0
	}
}

func (*NewRuntimeOptions) Inheritor(v Runtime) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.inheritor = v
	}
}

func (*NewRuntimeOptions) InitFunc(v func(rt Runtime)) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.initFunc = v
	}
}

func (*NewRuntimeOptions) StartFunc(v func(rt Runtime)) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.startFunc = v
	}
}

func (*NewRuntimeOptions) StopFunc(v func(rt Runtime)) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.stopFunc = v
	}
}

func (*NewRuntimeOptions) EnableAutoRun(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.enableAutoRun = v
	}
}

func (*NewRuntimeOptions) EnableAutoRecover(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.enableAutoRecover = v
	}
}

func (*NewRuntimeOptions) SafeCallCacheSize(v int) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.safeCallCacheSize = v
	}
}

func (*NewRuntimeOptions) FrameCreatorFunc(v func(rt Runtime) Frame) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.frameCreatorFunc = v
	}
}

func (*NewRuntimeOptions) EnableGC(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.enableGC = v
	}
}

func (*NewRuntimeOptions) GCTimeInterval(v time.Duration) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.gcTimeInterval = v
	}
}

func (*NewRuntimeOptions) GCItemNum(v int) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.gcItemNum = v
	}
}

func (*NewRuntimeOptions) Cache(v *misc.Cache) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.cache = v
	}
}

func (*NewRuntimeOptions) EnableEventRecursion(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.enableEventRecursion = v
	}
}

func (*NewRuntimeOptions) DiscardRecursiveEvent(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.discardRecursiveEvent = v
	}
}

func (*NewRuntimeOptions) CallEventDepth(v int) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.callEventDepth = v
	}
}
