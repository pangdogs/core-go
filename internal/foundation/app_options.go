package foundation

import "github.com/pangdogs/core/internal"

var NewAppOption = &NewAppOptions{}

type AppOptions struct {
	inheritor AppWhole
	initFunc,
	startFunc,
	stopFunc func()
	autoRecover bool
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
		o.inheritor = v.(AppWhole)
	}
}

func (*NewAppOptions) AutoRecover(v bool) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.autoRecover = v
	}
}

func (*NewAppOptions) InitFunc(v func()) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.initFunc = v
	}
}

func (*NewAppOptions) StartFunc(v func()) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.startFunc = v
	}
}

func (*NewAppOptions) StopFunc(v func()) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.stopFunc = v
	}
}
