package tasks

import (
	"fmt"
	"time"
)

const (
	None = iota
	High
	Low
)

type Priority int

type Task struct {
	Title             string
	Due               string
	Priority          Priority
	Active            bool
	Complete          bool
	TotalActiveTime   time.Duration
	LastActivatedTime time.Time
	Resources         []Resource
	Status            []Status
	Steps             []Step
}

func (t *Task) FilterValue() string {
	return t.Title
}

func (t *Task) Activate() {
	t.Active = true
	t.LastActivatedTime = time.Now()
}

func (t *Task) Deactivate() {
	t.Active = false
	t.TotalActiveTime += time.Now().Sub(t.LastActivatedTime)
}

// TODO this go to into the renders file...
func (t *Task) Runtime() string {
	seconds := int((time.Now().Add(t.TotalActiveTime).Sub(t.LastActivatedTime)).Seconds())
	switch {
	case seconds < 60:
		return fmt.Sprintf("%d sec", seconds)
	case seconds < 3600:
		return fmt.Sprintf("%d min %d sec", seconds/60, seconds%60)
	case seconds < 86400:
		return fmt.Sprintf("%d hr %d min %d sec", seconds/60/60, (seconds/60)%60, seconds%60)
	default:
		return fmt.Sprintf("%d days %d hr %d min %d sec", seconds/60/60/24, (seconds/60/60)%24, (seconds/60)%60, seconds%60)
	}
}

func ActiveTask(tasks []Task) int {
	for i, task := range tasks {
		if task.Active {
			return i
		}
	}
	return -1
}

func HasActiveTask(tasks []Task) bool {
	for _, task := range tasks {
		if task.Active {
			return true
		}
	}
	return false
}
