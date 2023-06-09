package timewheel

import "container/list"

func (this *TimeWheel) RunToTime(nowTime int64) {
	if this.startTime+this.curTime < nowTime {
		this.advanceScale()
	}
}

// advanceScale 前进一个刻度
func (this *TimeWheel) advanceScale() {
	for tmpOneWheel := this.firstWheel; tmpOneWheel != nil; tmpOneWheel = tmpOneWheel.nextWheel {
		if tmpOneWheel.level == 0 {
			this.doTask()
			this.curTime += tmpOneWheel.tickScale
			tmpOneWheel.advancePos()
		} else {
			tmpOneWheel.diffuseTask(this)
			tmpOneWheel.advancePos()
		}
		if tmpOneWheel.curPos != 0 {
			break
		}
	}
}

func (this *TimeWheel) doTask() {
	var l *list.List
	if l = this.firstWheel.slots[this.firstWheel.curPos]; l == nil {
		return
	}
	for elem := l.Front(); elem != nil; elem = elem.Next() {
		task, _ := elem.Value.(*oneTask)
		if task.callback != nil {
			// 直接goruntine执行好吗
			// go task.callback()
			task.callback()
		}
		if task.periodTime > 0 {
			task.targetTime += task.periodTime
			this.addTask(task)
		}
	}
	this.firstWheel.slots[this.firstWheel.curPos] = nil
}

func (this *oneWheel) advancePos() {
	this.curPos++
	this.curPos %= this.wheelSize
}

func (this *oneWheel) diffuseTask(wheel *TimeWheel) {
	// 直接重新调用 addTask 应该是可以的把
	var l *list.List
	if l = this.slots[this.curPos]; l == nil {
		return
	}
	for elem := l.Front(); elem != nil; elem = elem.Next() {
		task, _ := elem.Value.(*oneTask)
		wheel.addTask(task)
	}
	this.slots[this.curPos] = nil
}
