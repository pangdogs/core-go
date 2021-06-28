package foundation

type EntityInheritorWhole interface {
	initEntityInheritor(e EntityWhole)
}

type EntityInheritor struct {
	EntityWhole
}

func (ei *EntityInheritor) initEntityInheritor(e EntityWhole) {
	if e == nil {
		panic("nil entity")
	}
	ei.EntityWhole = e
}
