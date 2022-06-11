package container

type GC interface {
	GC()
	MarkGC()
	NeedGC() bool
}
