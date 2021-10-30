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
	autoRun           bool
	autoRecover       bool
	safeCallCacheSize int
	frameCreatorFunc  func(rt Runtime) Frame
	gcEnable          bool
	gcTimeInterval    time.Duration
	gcItemNum         int
	cache             *misc.Cache
}

type NewRuntimeOptionFunc func(o *RuntimeOptions)

type NewRuntimeOptions struct{}

func (*NewRuntimeOptions) Default() NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.inheritor = nil
		o.initFunc = nil
		o.startFunc = nil
		o.stopFunc = nil
		o.autoRun = false
		o.autoRecover = false
		o.safeCallCacheSize = 100
		o.frameCreatorFunc = nil
		o.gcEnable = true
		o.gcTimeInterval = 10 * time.Second
		o.gcItemNum = 1000
		o.cache = nil
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

func (*NewRuntimeOptions) AutoRun(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.autoRun = v
	}
}

func (*NewRuntimeOptions) AutoRecover(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.autoRecover = v
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

func (*NewRuntimeOptions) GCEnable(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.gcEnable = v
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
