package foundation

import "github.com/pangdogs/core/internal"

var NewEntityOption = &NewEntityOptions{}

type EntityOptions struct {
	inheritor EntityWhole
	initFunc,
	shutFunc func(entity internal.Entity)
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
		o.inheritor = v.(EntityWhole)
	}
}

func (*NewEntityOptions) InitFunc(v func(entity internal.Entity)) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.initFunc = v
	}
}

func (*NewEntityOptions) ShutFunc(v func(entity internal.Entity)) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.shutFunc = v
	}
}
