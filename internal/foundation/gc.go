package foundation

type GC interface {
	GC()
	GCHandle() uintptr
}

type GCRoot interface {
	PushGC(gc GC)
	RunGC()
}
