package core

import "time"

var NewRuntimeOption = &NewRuntimeOptions{}

type RuntimeOptions struct {
	Inheritor            Runtime
	EnableAutoRun        bool
	EnableAutoRecover    bool
	ProcessQueueCapacity int
	ProcessQueueTimeout  time.Duration
	Frame                Frame
}

type NewRuntimeOptionFunc func(o *RuntimeOptions)

type NewRuntimeOptions struct{}

func (*NewRuntimeOptions) Default() NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.Inheritor = nil
		o.EnableAutoRun = false
		o.EnableAutoRecover = false
		o.ProcessQueueCapacity = 128
		o.ProcessQueueTimeout = 5 * time.Second
		o.Frame = nil
	}
}

func (*NewRuntimeOptions) Inheritor(v Runtime) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.Inheritor = v
	}
}

func (*NewRuntimeOptions) EnableAutoRun(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.EnableAutoRun = v
	}
}

func (*NewRuntimeOptions) EnableAutoRecover(v bool) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.EnableAutoRecover = v
	}
}

func (*NewRuntimeOptions) ProcessQueueCapacity(v int) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.ProcessQueueCapacity = v
	}
}

func (*NewRuntimeOptions) ProcessQueueTimeout(v time.Duration) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.ProcessQueueTimeout = v
	}
}

func (*NewRuntimeOptions) Frame(v Frame) NewRuntimeOptionFunc {
	return func(o *RuntimeOptions) {
		o.Frame = v
	}
}
