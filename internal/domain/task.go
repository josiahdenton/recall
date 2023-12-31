package domain

// TODO rename to be more accurate
const (
	TaskPriorityNone = iota
	TaskPriorityLow
	TaskPriorityHigh
)

type Priority int

type Task struct {
	Title     string
	Due       string
	Priority  Priority
	Active    bool
	Complete  bool
	Resources []Resource
	Status    []Status
	Steps     []Step
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
