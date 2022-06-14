package container

type GC interface {
	GC() bool
	MarkGC()
	NeedGC() bool
}
