package repository

import "github.com/josiahdenton/recall/internal/domain"

type Repository interface {
	Task(string) *domain.Task
	Cycle(string) *domain.Cycle // etc...
	Accomplishment(string) *domain.Accomplishment
	SaveAccomplishment(domain.Accomplishment)
	AllAccomplishments([]string) []domain.Accomplishment
	SaveTask(domain.Task)
	AllTasks() []domain.Task
	SaveCycle(domain.Cycle)
	AllCycles() []domain.Cycle
	SaveChanges()
	// Add the following
	// - Zettel
	// - Release
	// - Resource
	// - Artifact
}
