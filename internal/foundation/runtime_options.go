package foundation

import "github.com/pangdogs/core/internal"

var NewRuntimeOption = &NewRuntimeOptions{}

type RuntimeOptions struct {
	inheritor RuntimeWhole
	initFunc,
	startFunc,
	stopFunc func(rt internal.Runtime)
	autoRun           bool
	autoRecover       bool
	safeCallCacheSize int
	frameCreatorFunc  func(rt internal.Runtime) internal.Frame
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
		o.inheritor = v.(RuntimeWhole)
	}
}

func (*NewRuntimeOptions) InitFunc(v func(rt internal.Runtime)) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.initFunc = v
	}
}

func (*NewRuntimeOptions) StartFunc(v func(rt internal.Runtime)) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.startFunc = v
	}
}

func (*NewRuntimeOptions) StopFunc(v func(rt internal.Runtime)) NewRuntimeOptionFunc {
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

func (*NewRuntimeOptions) FrameCreatorFunc(v func(rt internal.Runtime) internal.Frame) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.frameCreatorFunc = v
	}
}
