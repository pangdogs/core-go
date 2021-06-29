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
	CycleEnd()
	FrameBegin()
	FrameEnd()
	UpdateBegin()
	UpdateEnd()
}

func NewFrame(targetFPS float32, totalFrames uint64) internal.Frame {
	frame := &Frame{}
	frame.InitFrame(targetFPS, totalFrames)
	return frame
}

type Frame struct {
	targetFPS, curFPS      float32
	totalFrames, curFrames uint64
	curFrameBeginTime      time.Time
	lastFrameElapseTime    time.Duration
	curUpdateBeginTime     time.Time
	lastUpdateElapseTime   time.Duration
	statFPSBeginTime       time.Time
	statFPSFrames          uint64
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

func (f *Frame) GetCurFrameBeginTime() time.Time {
	return f.curFrameBeginTime
}

func (f *Frame) GetLastFrameElapseTime() time.Duration {
	return f.lastFrameElapseTime
}

func (f *Frame) GetCurUpdateBeginTime() time.Time {
	return f.curUpdateBeginTime
}

func (f *Frame) GetLastUpdateElapseTime() time.Duration {
	return f.lastUpdateElapseTime
}

func (f *Frame) SetCurFrames(v uint64) {
	f.curFrames = v
}

func (f *Frame) CycleBegin() {
	now := time.Now()

	f.curFPS = 0
	f.curFrames = 0

	f.statFPSBeginTime = now
	f.statFPSFrames = 0

	f.curFrameBeginTime = now
	f.lastFrameElapseTime = 0

	f.curUpdateBeginTime = now
	f.lastUpdateElapseTime = 0
}

func (f *Frame) CycleEnd() {
}

func (f *Frame) FrameBegin() {
	now := time.Now()

	statInterval := now.Sub(f.statFPSBeginTime).Seconds()
	if statInterval >= 1 {
		f.curFPS = float32(float64(f.statFPSFrames) / statInterval)
		f.statFPSBeginTime = now
		f.statFPSFrames = 0
	}

	f.curFrameBeginTime = now
}

func (f *Frame) FrameEnd() {
	f.lastFrameElapseTime = time.Now().Sub(f.curFrameBeginTime)
	f.statFPSFrames++
}

func (f *Frame) UpdateBegin() {
	f.curUpdateBeginTime = time.Now()
}

func (f *Frame) UpdateEnd() {
	f.lastUpdateElapseTime = time.Now().Sub(f.curUpdateBeginTime)
}
