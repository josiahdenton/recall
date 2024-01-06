package repository

import (
	"github.com/josiahdenton/recall/internal/domain"
	"time"
)

type LocalStorage struct {
	tasks  []domain.Task
	cycles []domain.Cycle
}

func NewLocalStorage() *LocalStorage {
	cycle := domain.NewCycle("EOY 2023", time.Now())
	cycle.Accomplishments = []domain.Accomplishment{
		domain.NewAccomplishment("Created a new cli tool", "it's epic", "Job Skills"),
		domain.NewAccomplishment("Became tech lead", "kinda scary", "Teamwork"),
	}
	cycle.Active = true

	return &LocalStorage{
		tasks: []domain.Task{
			domain.NewTask("wash dishes", "01/04/2024", domain.TaskPriorityLow),
			domain.NewTask("take out trash", "01/04/2024", domain.TaskPriorityLow),
			domain.NewTask("release ee decider", "01/05/2024", domain.TaskPriorityLow),
		},
		cycles: []domain.Cycle{
			cycle,
		},
	}
}

func (l *LocalStorage) Task(id string) domain.Task {
	for _, task := range l.tasks {
		if task.Id == id {
			return task
		}
	}
	return domain.Task{}
}

func (l *LocalStorage) SaveTask(task domain.Task) {
	// if id exists, don't just append... replace
	if l.taskExists(task.Id) {
		for i := range l.tasks {
			if l.tasks[i].Id == task.Id {
				l.tasks[i] = task
			}
		}
	} else {
		l.tasks = append(l.tasks, task)
	}
}

func (l *LocalStorage) taskExists(id string) bool {
	for _, task := range l.tasks {
		if task.Id == id {
			return true
		}
	}
	return false
}

func (l *LocalStorage) AllTasks() []domain.Task {
	return l.tasks
}

func (l *LocalStorage) Cycle(id string) domain.Cycle {
	for _, cycle := range l.cycles {
		if cycle.Id == id {
			return cycle
		}
	}
	return domain.Cycle{}
}

func (l *LocalStorage) SaveCycle(cycle domain.Cycle) {
	if l.cycleExists(cycle.Id) {
		for i := range l.cycles {
			if l.cycles[i].Id == cycle.Id {
				l.cycles[i] = cycle
			}
		}
	} else {
		l.cycles = append(l.cycles, cycle)
	}
}

func (l *LocalStorage) cycleExists(id string) bool {
	for _, cycle := range l.cycles {
		if cycle.Id == id {
			return true
		}
	}
	return false
}

func (l *LocalStorage) AllCycles() []domain.Cycle {
	return l.cycles
}

func (l *LocalStorage) Accomplishment(id string) domain.Accomplishment {
	for _, cycle := range l.cycles {
		for _, accomplishment := range cycle.Accomplishments {
			if accomplishment.Id == id {
				return accomplishment
			}
		}
	}
	return domain.Accomplishment{}
}

func (l *LocalStorage) SaveChanges() {
	return
}
