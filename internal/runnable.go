package internal

type Runnable interface {
	Run() chan struct{}
	Stop()
}
