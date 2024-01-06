package repository

import "github.com/josiahdenton/recall/internal/domain"

// TODO - internally to the repository, everything is only referred to each other by keys
// but when fetched from the repository, it will give everything back

type Repository interface {
	Task(string) domain.Task
	Cycle(string) domain.Cycle // etc...
	Accomplishment(string) domain.Accomplishment
	SaveTask(domain.Task)
	AllTasks() []domain.Task
	SaveCycle(domain.Cycle)
	AllCycles() []domain.Cycle
	SaveChanges()
}
