package container

type GC interface {
	GC() bool
	MarkGC()
	NeedGC() bool
}

type GCCollector interface {
	CollectGC(gc GC)
}
