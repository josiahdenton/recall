package domain

import (
	"gorm.io/gorm"
	"reflect"
	"time"
)

// Due TODO should be changed to a time.Time
// TODO this should be a list of Resource Ids
// TODO tasks can be recurring
// TODO tasks have a difficulty rating
// avg Task completion time is tracked

type Task struct {
	gorm.Model
	Title       string
	Tags        string
	Description string
	Due         time.Time
	// Active tracks if the current task is actively being worked on
	Active bool
	// Duration is the total time worked on a task
	Duration time.Duration
	// LastActivatedTime tracks how long it's been since this task was activated
	LastActivatedTime time.Time
	Archive           bool
	Resources         []Resource `gorm:"many2many:task_resources"`
	Steps             []Step
	AccomplishmentID  uint
}

func (t *Task) ToggleActive() {
	if t.Active && reflect.ValueOf(t.Duration).IsZero() {
		t.Duration = time.Now().Sub(t.LastActivatedTime)
	} else if t.Active {
		t.Duration += time.Now().Sub(t.LastActivatedTime)
	} else {
		t.LastActivatedTime = time.Now()
	}
	t.Active = !t.Active
}

// ActiveDuration calculates the total time
// this task has been actively worked on
func (t *Task) ActiveDuration() time.Duration {
	if t.Active && reflect.ValueOf(t.Duration).IsZero() {
		return time.Now().Sub(t.LastActivatedTime)
	} else if t.Active {
		return t.Duration + time.Now().Sub(t.LastActivatedTime)
	} else {
		return t.Duration
	}
}

func (t *Task) RemoveResource(i int) {
	t.Resources = append(t.Resources[:i], t.Resources[i+1:]...)
}

func (t *Task) RemoveStep(i int) {
	t.Steps = append(t.Steps[:i], t.Steps[i+1:]...)
}

func (t *Task) FilterValue() string {
	return t.Title + t.Tags
}
