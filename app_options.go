package core

var NewAppOption = &NewAppOptions{}

type AppOptions struct {
	Inheritor         App
	EnableAutoRecover bool
}

type NewAppOptionFunc func(o *AppOptions)

type NewAppOptions struct{}

func (*NewAppOptions) Default() NewAppOptionFunc {
	return func(o *AppOptions) {
		o.Inheritor = nil
		o.EnableAutoRecover = false
	}
}

func (*NewAppOptions) Inheritor(v App) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.Inheritor = v
	}
}

func (*NewAppOptions) EnableAutoRecover(v bool) NewAppOptionFunc {
	return func(o *AppOptions) {
		o.EnableAutoRecover = v
	}
}
