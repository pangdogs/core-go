package foundation

type GC interface {
	GC()
}

type GCRoot interface {
	PushGC(gc GC)
}
