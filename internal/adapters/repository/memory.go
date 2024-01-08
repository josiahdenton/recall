package repository

import (
	"github.com/josiahdenton/recall/internal/domain"
	"time"
)

type InMemoryStorage struct {
	tasks           []domain.Task
	cycles          []domain.Cycle
	accomplishments map[string]domain.Accomplishment
}

func NewInMemoryStorage() *InMemoryStorage {
	cycle := domain.NewCycle("EOY 2023", time.Now())

	acc1 := domain.NewAccomplishment("Created a new cli tool", "it's epic", "Job Skills")
	acc2 := domain.NewAccomplishment("Became tech lead", "kinda scary", "Teamwork")
	accomplishments := make(map[string]domain.Accomplishment)
	accomplishments[acc1.Id] = acc1
	accomplishments[acc2.Id] = acc2

	cycle.AccomplishmentIds = make([]string, 0)
	cycle.AccomplishmentIds = append(cycle.AccomplishmentIds, acc1.Id, acc2.Id)
	cycle.Active = true

	return &InMemoryStorage{
		tasks: []domain.Task{
			domain.NewTask("wash dishes", "01/04/2024", domain.TaskPriorityLow),
			domain.NewTask("take out trash", "01/04/2024", domain.TaskPriorityLow),
			domain.NewTask("release ee decider", "01/05/2024", domain.TaskPriorityLow),
		},
		cycles: []domain.Cycle{
			cycle,
		},
		accomplishments: accomplishments,
	}
}

func (l *InMemoryStorage) Task(id string) *domain.Task {
	for _, task := range l.tasks {
		if task.Id == id {
			return &task
		}
	}
	return &domain.Task{}
}

func (l *InMemoryStorage) SaveTask(task domain.Task) {
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

func (l *InMemoryStorage) taskExists(id string) bool {
	for _, task := range l.tasks {
		if task.Id == id {
			return true
		}
	}
	return false
}

func (l *InMemoryStorage) AllTasks() []domain.Task {
	return l.tasks
}

func (l *InMemoryStorage) Cycle(id string) *domain.Cycle {
	for _, cycle := range l.cycles {
		if cycle.Id == id {
			return &cycle
		}
	}
	return &domain.Cycle{}
}

func (l *InMemoryStorage) SaveCycle(cycle domain.Cycle) {
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

func (l *InMemoryStorage) cycleExists(id string) bool {
	for _, cycle := range l.cycles {
		if cycle.Id == id {
			return true
		}
	}
	return false
}

func (l *InMemoryStorage) AllCycles() []domain.Cycle {
	return l.cycles
}

func (l *InMemoryStorage) AllAccomplishments(ids []string) []domain.Accomplishment {
	accomplishments := make([]domain.Accomplishment, len(ids))
	for i, id := range ids {
		accomplishments[i] = l.accomplishments[id]
	}
	return accomplishments
}

func (l *InMemoryStorage) Accomplishment(id string) *domain.Accomplishment {
	accomplishment, ok := l.accomplishments[id]
	if !ok {
		return nil
	}
	return &accomplishment
}

func (l *InMemoryStorage) SaveAccomplishment(accomplishment domain.Accomplishment) {
	if _, ok := l.accomplishments[accomplishment.Id]; !ok {
		l.accomplishments[accomplishment.Id] = accomplishment
	}
}

func (l *InMemoryStorage) SaveChanges() {
	return
}

func (l *InMemoryStorage) SaveSettings(settings domain.Settings) {
	return
}
