package tasks

// TODO rename to be more accurate
const (
	None = iota
	Low
	High
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
