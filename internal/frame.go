package internal

import "time"

type Frame interface {
	GetTargetFPS() float32
	GetCurFPS() float32
	GetTotalFrames() uint64
	GetCurFrames() uint64
	GetCurFrameBeginTime() time.Time
	GetLastFrameElapseTime() time.Duration
	GetCurUpdateBeginTime() time.Time
	GetLastUpdateElapseTime() time.Duration
}
