package repository

import "github.com/josiahdenton/recall/internal/domain"

type Repository interface {
	Task(string) domain.Task
	Cycle(string) domain.Cycle // etc...
	SaveTask(domain.Task)
	AllTasks() []domain.Task
	SaveCycle(domain.Cycle)
	AllCycles() []domain.Cycle
	SaveChanges()
}
