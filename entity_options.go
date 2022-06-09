package core

import (
	"github.com/pangdogs/core/container"
)

var NewEntityOption = &NewEntityOptions{}

type EntityOptions struct {
	Inheritor                  Entity
	FaceCache                  *container.Cache[Face]
	HookCache                  *container.Cache[Hook]
	EnableFastGetComponent     bool
	EnableFastGetComponentByID bool
}

type NewEntityOptionFunc func(o *EntityOptions)

type NewEntityOptions struct{}

func (*NewEntityOptions) Default() NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.Inheritor = nil
		o.FaceCache = nil
		o.HookCache = nil
		o.EnableFastGetComponent = false
		o.EnableFastGetComponentByID = false
	}
}

func (*NewEntityOptions) Inheritor(v Entity) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.Inheritor = v
	}
}

func (*NewEntityOptions) FaceCache(v *container.Cache[Face]) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.FaceCache = v
	}
}

func (*NewEntityOptions) HookCache(v *container.Cache[Hook]) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.HookCache = v
	}
}

func (*NewEntityOptions) EnableFastGetComponent(v bool) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.EnableFastGetComponent = v
	}
}

func (*NewEntityOptions) EnableFastGetComponentByID(v bool) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.EnableFastGetComponentByID = v
	}
}
