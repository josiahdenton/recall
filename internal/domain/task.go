package domain

import (
	"gorm.io/gorm"
	"time"
)

// TODO rename to be more accurate
const (
	TaskPriorityNone = iota
	TaskPriorityLow
	TaskPriorityHigh
)

type Priority int

// Due TODO should be changed to a time.Time
// TODO this should be a list of Resource Ids
// TODO tasks can be recurring
// TODO tasks have a difficulty rating
// avg Task completion time is tracked

type Task struct {
	gorm.Model
	Title            string
	Due              time.Time
	Priority         Priority
	Active           bool
	Archive          bool
	Favorite         bool
	Resources        []Resource
	Status           []Status
	Steps            []Step
	AccomplishmentID uint
}

func (t *Task) RemoveResource(i int) {
	t.Resources = append(t.Resources[:i], t.Resources[i+1:]...)
}

func (t *Task) RemoveStatus(i int) {
	t.Status = append(t.Status[:i], t.Status[i+1:]...)
}

func (t *Task) RemoveStep(i int) {
	t.Steps = append(t.Steps[:i], t.Steps[i+1:]...)
}

func (t *Task) FilterValue() string {
	return t.Title
}
