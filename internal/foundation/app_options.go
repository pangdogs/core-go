package foundation

import "github.com/pangdogs/core/internal"

var NewAppOption = &NewAppOptions{}

type AppOptions struct {
	inheritor internal.App
}

type NewAppOptionFunc func(o *AppOptions)

type NewAppOptions struct{}

func (*NewAppOptions) Default() NewAppOptionFunc {
	return func(o *AppOptions) {
		o.inheritor = nil
	}
}

func (*NewAppOptions) Inheritor(v internal.App) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.inheritor = v
	}
}
