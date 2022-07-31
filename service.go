package core

type Service interface {
	Runnable
	init(ctx ServiceContext, opts *ServiceOptions)
	getOptions() *ServiceOptions
	GetContext() ServiceContext
}

func ServiceGetOptions(serv Service) ServiceOptions {
	return *serv.getOptions()
}

func ServiceGetInheritor(serv Service) Service {
	return serv.getOptions().Inheritor
}

func NewService(servCtx ServiceContext, optFuncs ...NewServiceOptionFunc) Service {
	opts := &ServiceOptions{}
	NewServiceOption.Default()(opts)

	for i := range optFuncs {
		optFuncs[i](opts)
	}

	if opts.Inheritor != nil {
		opts.Inheritor.init(servCtx, opts)
		return opts.Inheritor
	}

	serv := &ServiceBehavior{}
	serv.init(servCtx, opts)

	return serv.opts.Inheritor
}

type ServiceBehavior struct {
	opts ServiceOptions
	ctx  ServiceContext
}

func (serv *ServiceBehavior) init(servCtx ServiceContext, opts *ServiceOptions) {
	if servCtx == nil {
		panic("nil servCtx")
	}

	if opts == nil {
		panic("nil opts")
	}

	serv.opts = *opts

	if serv.opts.Inheritor == nil {
		serv.opts.Inheritor = serv
	}

	serv.ctx = servCtx
}

func (serv *ServiceBehavior) getOptions() *ServiceOptions {
	return &serv.opts
}

func (serv *ServiceBehavior) GetContext() ServiceContext {
	return serv.ctx
}
