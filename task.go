package timewheel

import (
	"container/list"
	"fmt"
)

func (this *TimeWheel) addTask(task *oneTask) (int64, error) {
	diffTime := this.fixDiffTime(task)
	for tmpOneWheel := this.firstWheel; tmpOneWheel != nil; tmpOneWheel = tmpOneWheel.nextWheel {
		// 可以放进当前时间轮
		if diffTime < tmpOneWheel.getResidueTime() {
			putPos := diffTime/tmpOneWheel.tickScale + tmpOneWheel.curPos
			if tmpOneWheel.slots[putPos] == nil {
				tmpOneWheel.slots[putPos] = &list.List{}
			}
			tmpOneWheel.slots[putPos].PushBack(task)
			return task.guid, nil
		}
	}
	// 所有时间轮都放不下,后续可能需要增加错误处理
	return 0, fmt.Errorf("not push task, targetTime %d", task.targetTime)
}

func (this *TimeWheel) AddTask(targetTime, periodTime int64, callback TimeWheelCallBack) {
	this.genTaskGUID++
	this.addTask(&oneTask{guid: this.genTaskGUID, targetTime: targetTime, periodTime: periodTime, callback: callback})
}

// RemoveTask 移除guid标识的任务, 只加入标记，在使用的时候跳过处理
// 这样好像不是特别好，如果是有频繁的对未来的任务增加和移除, 会导致 guidRemoveFlag 一直累加
func (this *TimeWheel) RemoveTask(guid int64) {
	this.guidRemoveFlag[guid] = struct{}{}
}

func (this *TimeWheel) fixDiffTime(task *oneTask) int64 {
	diffTime := task.targetTime - this.startTime - this.curTime
	if diffTime < 0 {
		// 对于周期任务需要, 不能立马执行，需要修正
		// 比如周期任务 &oneTask{targetTime: 1, periodTime: 60}, 当前时间是30
		// 如果直接返回0的话, 会导致后续执行的是  30 90 150..., 但是想要的应该是 1 61 121...
		// 不过感觉，是不是还是要让他立马执行一次阿
		if task.periodTime > 0 {
			return diffTime%task.periodTime + task.periodTime
		}
		// 非周期任务直接返回0, 即立马执行
		return 0
	}
	return diffTime
}
