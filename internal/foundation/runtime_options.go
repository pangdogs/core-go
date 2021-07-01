package foundation

import "github.com/pangdogs/core/internal"

var NewRuntimeOption = &NewRuntimeOptions{}

type RuntimeOptions struct {
	inheritor internal.Runtime
	initFunc,
	startFunc,
	stopFunc func()
	autoRun           bool
	autoRecover       bool
	reportError       chan error
	safeCallCacheSize int
	frameCreatorFunc  func() internal.Frame
}

type NewRuntimeOptionFunc func(o *RuntimeOptions)

type NewRuntimeOptions struct{}

func (*NewRuntimeOptions) Default() NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.autoRun = false
		o.autoRecover = true
		o.safeCallCacheSize = 100
		o.frameCreatorFunc = nil
	}
}

func (*NewRuntimeOptions) Inheritor(v internal.Runtime) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.inheritor = v
	}
}

func (*NewRuntimeOptions) InitFunc(v func()) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.initFunc = v
	}
}

func (*NewRuntimeOptions) StartFunc(v func()) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.startFunc = v
	}
}

func (*NewRuntimeOptions) StopFunc(v func()) NewRuntimeOptionFunc {
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

func (*NewRuntimeOptions) ReportError(v chan error) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.reportError = v
	}
}

func (*NewRuntimeOptions) SafeCallCacheSize(v int) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.safeCallCacheSize = v
	}
}

func (*NewRuntimeOptions) FrameCreatorFunc(v func() internal.Frame) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.frameCreatorFunc = v
	}
}
