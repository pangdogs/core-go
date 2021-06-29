package foundation

type AppInheritorWhole interface {
	initAppInheritor(app AppWhole)
}

type AppInheritor struct {
	AppWhole
}

func (ai *AppInheritor) initAppInheritor(app AppWhole) {
	if app == nil {
		panic("nil app")
	}
	ai.AppWhole = app
}
