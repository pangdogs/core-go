package foundation

import "github.com/pangdogs/core/internal"

var NewEntityOption = &NewEntityOptions{}

type EntityOptions struct {
	inheritor internal.Entity
	initFunc,
	shutFunc func()
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

func (*NewEntityOptions) ShutFunc(v func()) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.shutFunc = v
	}
}
