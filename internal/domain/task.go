package domain

import (
	"github.com/google/uuid"
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
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Due       time.Time  `json:"due"`
	Priority  Priority   `json:"priority"`
	Active    bool       `json:"active"`
	Archive   bool       `json:"archive"`
	Resources []Resource `json:"resources"`
	Status    []Status   `json:"status"`
	Steps     []Step     `json:"steps"`
}

func NewTask(title string, due time.Time, priority Priority) Task {
	id, err := uuid.NewRandom()
	if err != nil {
		return Task{}
	}

	return Task{
		Id:       id.String(),
		Title:    title,
		Due:      due,
		Priority: priority,
	}
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
