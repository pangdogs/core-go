package foundation

import (
	"time"
)

type Frame interface {
	GetTargetFPS() float32
	GetCurFPS() float32
	GetTotalFrames() uint64
	GetCurFrames() uint64
	IsBlink() bool
	GetCycleBeginTime() time.Time
	GetCycleElapseTime() time.Duration
	GetCurFrameBeginTime() time.Time
	GetLastFrameElapseTime() time.Duration
	GetCurUpdateBeginTime() time.Time
	GetLastUpdateElapseTime() time.Duration
	setCurFrames(v uint64)
	cycleBegin()
	cycleEnd()
	frameBegin()
	frameEnd()
	updateBegin()
	updateEnd()
}

func NewFrame(targetFPS float32, totalFrames uint64, blink bool) Frame {
	frame := &_Frame{}
	frame.initFrame(targetFPS, totalFrames, blink)
	return frame
}

type _Frame struct {
	targetFPS, curFPS      float32
	totalFrames, curFrames uint64
	blink                  bool
	blinkFrameTime         time.Duration
	cycleBeginTime         time.Time
	cycleElapseTime        time.Duration
	curFrameBeginTime      time.Time
	lastFrameElapseTime    time.Duration
	curUpdateBeginTime     time.Time
	lastUpdateElapseTime   time.Duration
	statFPSBeginTime       time.Time
	statFPSFrames          uint64
}

func (f *_Frame) initFrame(targetFPS float32, totalFrames uint64, blink bool) {
	if targetFPS <= 0 {
		panic("[targetFPS > 0] is required")
	}

	if totalFrames < 0 {
		panic("[totalFrames >= 0] is required")
	}

	f.targetFPS = targetFPS
	f.totalFrames = totalFrames
	f.blink = blink

	if blink {
		f.blinkFrameTime = time.Duration(float64(time.Second) / float64(targetFPS))
	}
}

func (f *_Frame) GetTargetFPS() float32 {
	return f.targetFPS
}

func (f *_Frame) GetCurFPS() float32 {
	return f.curFPS
}

func (f *_Frame) GetTotalFrames() uint64 {
	return f.totalFrames
}

func (f *_Frame) GetCurFrames() uint64 {
	return f.curFrames
}

func (f *_Frame) IsBlink() bool {
	return f.blink
}

func (f *_Frame) GetCycleBeginTime() time.Time {
	return f.cycleBeginTime
}

func (f *_Frame) GetCycleElapseTime() time.Duration {
	return f.cycleElapseTime
}

func (f *_Frame) GetCurFrameBeginTime() time.Time {
	return f.curFrameBeginTime
}

func (f *_Frame) GetLastFrameElapseTime() time.Duration {
	return f.lastFrameElapseTime
}

func (f *_Frame) GetCurUpdateBeginTime() time.Time {
	return f.curUpdateBeginTime
}

func (f *_Frame) GetLastUpdateElapseTime() time.Duration {
	return f.lastUpdateElapseTime
}

func (f *_Frame) setCurFrames(v uint64) {
	f.curFrames = v
}

func (f *_Frame) cycleBegin() {
	now := time.Now()

	f.curFPS = 0

	f.statFPSBeginTime = now
	f.statFPSFrames = 0

	f.cycleBeginTime = now
	f.cycleElapseTime = 0

	f.curFrameBeginTime = now
	f.lastFrameElapseTime = 0

	f.curUpdateBeginTime = now
	f.lastUpdateElapseTime = 0
}

func (f *_Frame) cycleEnd() {
	if f.blink {
		f.curFPS = float32(float64(f.curFrames) / time.Now().Sub(f.cycleBeginTime).Seconds())
	}
}

func (f *_Frame) frameBegin() {
	now := time.Now()

	if !f.blink {
		statInterval := now.Sub(f.statFPSBeginTime).Seconds()
		if statInterval >= 1 {
			f.curFPS = float32(float64(f.statFPSFrames) / statInterval)
			f.statFPSBeginTime = now
			f.statFPSFrames = 0
		}
	}

	f.curFrameBeginTime = now
}

func (f *_Frame) frameEnd() {
	now := time.Now()

	if f.blink {
		f.cycleElapseTime += f.blinkFrameTime
	} else {
		f.cycleElapseTime = now.Sub(f.curFrameBeginTime)
	}

	f.lastFrameElapseTime = now.Sub(f.curFrameBeginTime)
	f.statFPSFrames++
}

func (f *_Frame) updateBegin() {
	f.curUpdateBeginTime = time.Now()
}

func (f *_Frame) updateEnd() {
	f.lastUpdateElapseTime = time.Now().Sub(f.curUpdateBeginTime)
}
