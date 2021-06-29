package foundation

type RuntimeInheritorWhole interface {
	initRuntimeInheritor(Runtime RuntimeWhole)
}

type RuntimeInheritor struct {
	RuntimeWhole
}

func (rti *RuntimeInheritor) initRuntimeInheritor(rt RuntimeWhole) {
	if rt == nil {
		panic("nil runtime")
	}
	rti.RuntimeWhole = rt
}
