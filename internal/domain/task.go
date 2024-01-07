package domain

import "github.com/google/uuid"

// TODO rename to be more accurate
const (
	TaskPriorityNone = iota
	TaskPriorityLow
	TaskPriorityHigh
)

type Priority int

type Task struct {
	Id    string
	Title string
	// Due TODO should be changed to a time.Time
	Due       string
	Priority  Priority
	Active    bool
	Complete  bool
	Resources []Resource // TODO this should be a list of Resource Ids
	Status    []Status
	Steps     []Step
}

func NewTask(title, due string, priority Priority) Task {
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
