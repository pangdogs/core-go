package core

var NewServiceContextOption = &NewServiceContextOptions{}

type ServiceContextOptions struct {
	Inheritor   ServiceContext
	ReportError chan error
	StartedCallback,
	StoppingCallback,
	StoppedCallback func(serv Service)
}

type NewServiceContextOptionFunc func(o *ServiceContextOptions)

type NewServiceContextOptions struct{}

func (*NewServiceContextOptions) Default() NewServiceContextOptionFunc {
	return func(o *ServiceContextOptions) {
		o.Inheritor = nil
		o.ReportError = nil
		o.StartedCallback = nil
		o.StoppingCallback = nil
		o.StoppedCallback = nil
	}
}

func (*NewServiceContextOptions) Inheritor(v ServiceContext) NewServiceContextOptionFunc {
	return func(o *ServiceContextOptions) {
		o.Inheritor = v
	}
}

func (*NewServiceContextOptions) ReportError(v chan error) NewServiceContextOptionFunc {
	return func(o *ServiceContextOptions) {
		o.ReportError = v
	}
}

func (*NewServiceContextOptions) StartedCallback(v func(serv Service)) NewServiceContextOptionFunc {
	return func(o *ServiceContextOptions) {
		o.StartedCallback = v
	}
}

func (*NewServiceContextOptions) StoppingCallback(v func(serv Service)) NewServiceContextOptionFunc {
	return func(o *ServiceContextOptions) {
		o.StoppingCallback = v
	}
}

func (*NewServiceContextOptions) StoppedCallback(v func(serv Service)) NewServiceContextOptionFunc {
	return func(o *ServiceContextOptions) {
		o.StoppedCallback = v
	}
}
