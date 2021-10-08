package foundation

var NewEntityOption = &NewEntityOptions{}

type EntityOptions struct {
	inheritor Entity
	initFunc,
	shutFunc func(entity Entity)
}

type NewEntityOptionFunc func(o *EntityOptions)

type NewEntityOptions struct{}

func (*NewEntityOptions) Default() NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.inheritor = nil
		o.initFunc = nil
		o.shutFunc = nil
	}
}

func (*NewEntityOptions) Inheritor(v Entity) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.inheritor = v
	}
}

func (*NewEntityOptions) InitFunc(v func(entity Entity)) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.initFunc = v
	}
}

func (*NewEntityOptions) ShutFunc(v func(entity Entity)) NewEntityOptionFunc {
	return func(o *EntityOptions) {
		o.shutFunc = v
	}
}
