package foundation

var NewAppOption = &NewAppOptions{}

type AppOptions struct {
	inheritor App
	initFunc,
	startFunc,
	stopFunc func(app App)
	autoRecover bool
}

type NewAppOptionFunc func(o *AppOptions)

type NewAppOptions struct{}

func (*NewAppOptions) Default() NewAppOptionFunc {
	return func(o *AppOptions) {
		o.inheritor = nil
		o.initFunc = nil
		o.startFunc = nil
		o.stopFunc = nil
		o.autoRecover = false
	}
}

func (*NewAppOptions) Inheritor(v App) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.inheritor = v
	}
}

func (*NewAppOptions) AutoRecover(v bool) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.autoRecover = v
	}
}

func (*NewAppOptions) InitFunc(v func(app App)) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.initFunc = v
	}
}

func (*NewAppOptions) StartFunc(v func(app App)) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.startFunc = v
	}
}

func (*NewAppOptions) StopFunc(v func(app App)) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.stopFunc = v
	}
}
