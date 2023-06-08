package timewheel

func (self *TimeWheel) addTask(task *oneTask) {
	relativeTime := task.targetTime - self.startTime
	for tmpOneWheel := self.firstWheel; tmpOneWheel != nil; tmpOneWheel = tmpOneWheel.nextWheel {

	}
}
