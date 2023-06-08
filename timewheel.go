package timewheel

import "container/list"

type TimeWheelCallBack func()

type oneTask struct {
	targetTime int64 // 目标时间
	periodTime int64 // 周期
	callback   TimeWheelCallBack
}

type oneWheel struct {
	tickScale int64 // 时间刻度
	wheelSize int64 // 时间轮的刻度数量
	level     int   // 当前时间轮层次
	prevWheel *oneWheel
	nextWheel *oneWheel
	slots     []*list.List
	curPos    int64
}

type TimeWheel struct {
	curTick    int64 // 当前时间
	startTime  int64 // 开始时间
	firstWheel *oneWheel
}

// NewTimeWheel 新建一个时间轮的结构, startTime 开始时间, firstScale 第一层时间轮一个刻度大小, allWheelSize 每层时间轮刻度数量
func NewTimeWheel(startTime, firstScale int64, allWheelSize []int64) *TimeWheel {
	if firstScale <= 0 || len(allWheelSize) == 0 {
		panic("Invalid parameter: firstScale allWheelSize")
	}
	timeWheel := &TimeWheel{startTime: startTime}
	var tmpOneWheel *oneWheel
	for i, wheelSize := range allWheelSize {
		if i == 0 {
			timeWheel.firstWheel = newOneWheel(wheelSize, nil)
			timeWheel.firstWheel.tickScale = firstScale
			tmpOneWheel = timeWheel.firstWheel
		} else {
			tmpOneWheel = newOneWheel(wheelSize, tmpOneWheel)
		}
	}
	return timeWheel
}

func newOneWheel(wheelSize int64, prevWheel *oneWheel) *oneWheel {
	newWheel := &oneWheel{}
	newWheel.wheelSize = wheelSize
	if prevWheel == nil {
		newWheel.level = 0
	} else {
		prevWheel.nextWheel = newWheel
		newWheel.prevWheel = prevWheel
		newWheel.level = prevWheel.level + 1
		newWheel.tickScale = prevWheel.tickScale * prevWheel.wheelSize
	}
	newWheel.slots = make([]*list.List, wheelSize)
	newWheel.curPos = 0
	return newWheel
}