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
			tmpOneWheel.advancePos()
			tmpOneWheel.diffuseTask(this)
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
		if this.isRemoved(task.guid) {
			continue
		}
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

// diffuseTask 将高层时间轮扩散到低层
func (this *oneWheel) diffuseTask(wheel *TimeWheel) {
	// 直接重新调用 addTask 应该是可以的把
	var l *list.List
	if l = this.slots[this.curPos]; l == nil {
		return
	}
	for elem := l.Front(); elem != nil; elem = elem.Next() {
		task, _ := elem.Value.(*oneTask)
		if wheel.isRemoved(task.guid) {
			continue
		}
        // TODO: 这里直接调用addTask
		wheel.addTask(task)
	}
	this.slots[this.curPos] = nil
}
