package core

import (
	"github.com/pangdogs/core/container"
)

var NewRuntimeContextOption = &NewRuntimeContextOptions{}

type RuntimeContextOptions struct {
	Inheritor RuntimeContext
	StartedCallback,
	StoppedCallback func(runtime Runtime)
	FaceCache *container.Cache[Face]
	HookCache *container.Cache[Hook]
}

type NewRuntimeContextOptionFunc func(o *RuntimeContextOptions)

type NewRuntimeContextOptions struct{}

func (*NewRuntimeContextOptions) Default() NewRuntimeContextOptionFunc {
	return func(o *RuntimeContextOptions) {
		o.Inheritor = nil
		o.StartedCallback = nil
		o.StoppedCallback = nil
		o.FaceCache = nil
		o.HookCache = nil
	}
}

func (*NewRuntimeContextOptions) Inheritor(v RuntimeContext) NewRuntimeContextOptionFunc {
	return func(o *RuntimeContextOptions) {
		o.Inheritor = v
	}
}

func (*NewRuntimeContextOptions) StartFunc(v func(rt Runtime)) NewRuntimeContextOptionFunc {
	return func(o *RuntimeContextOptions) {
		o.StartedCallback = v
	}
}

func (*NewRuntimeContextOptions) StopFunc(v func(rt Runtime)) NewRuntimeContextOptionFunc {
	return func(o *RuntimeContextOptions) {
		o.StoppedCallback = v
	}
}

func (*NewRuntimeContextOptions) FaceCache(v *container.Cache[Face]) NewRuntimeContextOptionFunc {
	return func(o *RuntimeContextOptions) {
		o.FaceCache = v
	}
}

func (*NewRuntimeContextOptions) HookCache(v *container.Cache[Hook]) NewRuntimeContextOptionFunc {
	return func(o *RuntimeContextOptions) {
		o.HookCache = v
	}
}
