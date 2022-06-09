package core

import (
	"time"
)

type Frame interface {
	GetTargetFPS() float32
	GetCurFPS() float32
	GetTotalFrames() uint64
	GetCurFrames() uint64
	setCurFrames(v uint64)
	Blink() bool
	GetRunningBeginTime() time.Time
	GetRunningElapseTime() time.Duration
	GetFrameBeginTime() time.Time
	GetLastFrameElapseTime() time.Duration
	GetUpdateBeginTime() time.Time
	GetLastUpdateElapseTime() time.Duration
	runningBegin()
	runningEnd()
	frameBegin()
	frameEnd()
	updateBegin()
	updateEnd()
}

func NewFrame(targetFPS float32, totalFrames uint64, blink bool) Frame {
	frame := &FrameBehavior{}
	frame.init(targetFPS, totalFrames, blink)
	return frame
}

type FrameBehavior struct {
	targetFPS, curFPS      float32
	totalFrames, curFrames uint64
	blink                  bool
	blinkFrameTime         time.Duration
	runningBeginTime       time.Time
	runningElapseTime      time.Duration
	frameBeginTime         time.Time
	lastFrameElapseTime    time.Duration
	updateBeginTime        time.Time
	lastUpdateElapseTime   time.Duration
	statFPSBeginTime       time.Time
	statFPSFrames          uint64
}

func (frame *FrameBehavior) init(targetFPS float32, totalFrames uint64, blink bool) {
	if targetFPS <= 0 {
		panic("targetFPS less equal 0 invalid")
	}

	if totalFrames < 0 {
		panic("totalFrames less 0 invalid")
	}

	frame.targetFPS = targetFPS
	frame.totalFrames = totalFrames
	frame.blink = blink

	if blink {
		frame.blinkFrameTime = time.Duration(float64(time.Second) / float64(targetFPS))
	}
}

func (frame *FrameBehavior) GetTargetFPS() float32 {
	return frame.targetFPS
}

func (frame *FrameBehavior) GetCurFPS() float32 {
	return frame.curFPS
}

func (frame *FrameBehavior) GetTotalFrames() uint64 {
	return frame.totalFrames
}

func (frame *FrameBehavior) GetCurFrames() uint64 {
	return frame.curFrames
}

func (frame *FrameBehavior) setCurFrames(v uint64) {
	frame.curFrames = v
}

func (frame *FrameBehavior) Blink() bool {
	return frame.blink
}

func (frame *FrameBehavior) GetRunningBeginTime() time.Time {
	return frame.runningBeginTime
}

func (frame *FrameBehavior) GetRunningElapseTime() time.Duration {
	return frame.runningElapseTime
}

func (frame *FrameBehavior) GetFrameBeginTime() time.Time {
	return frame.frameBeginTime
}

func (frame *FrameBehavior) GetLastFrameElapseTime() time.Duration {
	return frame.lastFrameElapseTime
}

func (frame *FrameBehavior) GetUpdateBeginTime() time.Time {
	return frame.updateBeginTime
}

func (frame *FrameBehavior) GetLastUpdateElapseTime() time.Duration {
	return frame.lastUpdateElapseTime
}

func (frame *FrameBehavior) runningBegin() {
	now := time.Now()

	frame.curFPS = 0
	frame.curFrames = 0

	frame.statFPSBeginTime = now
	frame.statFPSFrames = 0

	frame.runningBeginTime = now
	frame.runningElapseTime = 0

	frame.frameBeginTime = now
	frame.lastFrameElapseTime = 0

	frame.updateBeginTime = now
	frame.lastUpdateElapseTime = 0
}

func (frame *FrameBehavior) runningEnd() {
	if frame.blink {
		frame.curFPS = float32(float64(frame.curFrames) / time.Now().Sub(frame.runningBeginTime).Seconds())
	}
}

func (frame *FrameBehavior) frameBegin() {
	now := time.Now()

	frame.frameBeginTime = now

	if !frame.blink {
		statInterval := now.Sub(frame.statFPSBeginTime).Seconds()
		if statInterval >= 1 {
			frame.curFPS = float32(float64(frame.statFPSFrames) / statInterval)
			frame.statFPSBeginTime = now
			frame.statFPSFrames = 0
		}
	}
}

func (frame *FrameBehavior) frameEnd() {
	if frame.blink {
		frame.runningElapseTime += frame.blinkFrameTime
	} else {
		frame.lastFrameElapseTime = time.Now().Sub(frame.frameBeginTime)
		frame.runningElapseTime += frame.lastFrameElapseTime
		frame.statFPSFrames++
	}
}

func (frame *FrameBehavior) updateBegin() {
	frame.updateBeginTime = time.Now()
}

func (frame *FrameBehavior) updateEnd() {
	frame.lastUpdateElapseTime = time.Now().Sub(frame.updateBeginTime)
}
