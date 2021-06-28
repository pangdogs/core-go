package foundation

import (
	"github.com/pangdogs/core/internal"
	"time"
)

type FrameWhole interface {
	internal.Frame
	InitFrame(targetFPS float32, totalFrames uint64)
	SetCurFrames(v uint64)
	CycleBegin()
	FrameBegin()
	FrameEnd()
	CycleEnd()
}

func NewFrame(targetFPS float32, totalFrames uint64) internal.Frame {
	frame := &Frame{}
	frame.InitFrame(targetFPS, totalFrames)
	return frame
}

type Frame struct {
	targetFPS, curFPS      float32
	totalFrames, curFrames uint64
	curBeginTime           time.Time
	lastElapseTime         time.Duration
	statBeginTime          time.Time
	statFrames             uint64
}

func (f *Frame) InitFrame(targetFPS float32, totalFrames uint64) {
	f.targetFPS = targetFPS
	f.totalFrames = totalFrames
}

func (f *Frame) GetTargetFPS() float32 {
	return f.targetFPS
}

func (f *Frame) GetCurFPS() float32 {
	return f.curFPS
}

func (f *Frame) GetTotalFrames() uint64 {
	return f.totalFrames
}

func (f *Frame) GetCurFrames() uint64 {
	return f.curFrames
}

func (f *Frame) GetCurBeginTime() time.Time {
	return f.curBeginTime
}

func (f *Frame) GetLastElapseTime() time.Duration {
	return f.lastElapseTime
}

func (f *Frame) SetCurFrames(v uint64) {
	f.curFrames = v
}

func (f *Frame) CycleBegin() {
	f.statBeginTime = time.Now()
	f.statFrames = 0
	f.curFPS = 0
	f.curFrames = 0
	f.curBeginTime = f.statBeginTime
	f.lastElapseTime = 0
}

func (f *Frame) FrameBegin() {
	now := time.Now()

	statInterval := now.Sub(f.statBeginTime).Seconds()
	if statInterval >= 1 {
		f.curFPS = float32(float64(f.statFrames) / statInterval)
		f.statBeginTime = now
		f.statFrames = 0
	}

	f.curBeginTime = now
}

func (f *Frame) FrameEnd() {
	f.lastElapseTime = time.Now().Sub(f.curBeginTime)
	f.statFrames++
}

func (f *Frame) CycleEnd() {
}
