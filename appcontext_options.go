package core

var NewAppContextOption = &NewAppContextOptions{}

type AppContextOptions struct {
	Inheritor   AppContext
	ReportError chan error
	StartedCallback,
	StoppingCallback,
	StoppedCallback func(app App)
}

type NewAppContextOptionFunc func(o *AppContextOptions)

type NewAppContextOptions struct{}

func (*NewAppContextOptions) Default() NewAppContextOptionFunc {
	return func(o *AppContextOptions) {
		o.Inheritor = nil
		o.ReportError = nil
		o.StartedCallback = nil
		o.StoppingCallback = nil
		o.StoppedCallback = nil
	}
}

func (*NewAppContextOptions) Inheritor(v AppContext) NewAppContextOptionFunc {
	return func(o *AppContextOptions) {
		o.Inheritor = v
	}
}

func (*NewAppContextOptions) ReportError(v chan error) NewAppContextOptionFunc {
	return func(o *AppContextOptions) {
		o.ReportError = v
	}
}

func (*NewAppContextOptions) StartedCallback(v func(app App)) NewAppContextOptionFunc {
	return func(o *AppContextOptions) {
		o.StartedCallback = v
	}
}

func (*NewAppContextOptions) StoppingCallback(v func(app App)) NewAppContextOptionFunc {
	return func(o *AppContextOptions) {
		o.StoppingCallback = v
	}
}

func (*NewAppContextOptions) StoppedCallback(v func(app App)) NewAppContextOptionFunc {
	return func(o *AppContextOptions) {
		o.StoppedCallback = v
	}
}
