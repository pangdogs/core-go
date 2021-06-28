package foundation

import "github.com/pangdogs/core/internal"

var NewRuntimeOption = &NewRuntimeOptions{}

type RuntimeOptions struct {
	autoRun           bool
	autoRecover       bool
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

func (*NewRuntimeOptions) FrameCreatorFunc(v func() internal.Frame) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.frameCreatorFunc = v
	}
}
