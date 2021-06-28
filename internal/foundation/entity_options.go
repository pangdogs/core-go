package foundation

import "github.com/pangdogs/core/internal"

var NewEntityOption = &NewEntityOptions{}

type EntityOptions struct {
	inheritor internal.Entity
	initFunc,
	updateFunc,
	lateUpdateFunc,
	destroyFunc func()
}

type NewEntityOptionFunc func(o *EntityOptions)

type NewEntityOptions struct{}

func (*NewEntityOptions) Default() NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.inheritor = nil
	}
}

func (*NewEntityOptions) Inheritor(v internal.Entity) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.inheritor = v
	}
}

func (*NewEntityOptions) InitFunc(v func()) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.initFunc = v
	}
}

func (*NewEntityOptions) UpdateFunc(v func()) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.updateFunc = v
	}
}

func (*NewEntityOptions) LateUpdateFunc(v func()) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.lateUpdateFunc = v
	}
}

func (*NewEntityOptions) DestroyFunc(v func()) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.destroyFunc = v
	}
}
