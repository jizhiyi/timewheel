package timewheel

import (
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {
	start := time.Now().Unix()
	tw := NewTimeWheel(start, 1, []int64{60, 60, 24})
	t.Log(tw)
}

func TestRunToTime(t *testing.T) {
	start := int64(0)
	tw := NewTimeWheel(start, 1, []int64{60, 60, 24})
	// tw.AddTask(start, 2, func() {
	// 	t.Log("Task 1")
	// })
	// tw.AddTask(start, 3, func() {
	// 	t.Log("Task 2")
	// })
	tw.AddTask(start, 90, func() {
		t.Log("Task 3")
	})
	for i := 1; i < 300; i++ {
		t.Log("Cur", i-1)
		tw.RunToTime(start + int64(i))
	}
}
