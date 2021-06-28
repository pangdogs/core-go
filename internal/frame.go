package internal

import "time"

type Frame interface {
	GetTargetFPS() float32
	GetCurFPS() float32
	GetTotalFrames() uint64
	GetCurFrames() uint64
	GetCurBeginTime() time.Time
	GetLastElapseTime() time.Duration
}
