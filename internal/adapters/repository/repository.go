package repository

import (
	"github.com/josiahdenton/recall/internal/domain"
)

type Repository interface {
	Task(string) *domain.Task
	Cycle(string) *domain.Cycle // etc...
	Accomplishment(string) *domain.Accomplishment
	SaveAccomplishment(domain.Accomplishment)
	AllAccomplishments([]string) []domain.Accomplishment
	SaveTask(domain.Task)
	AllTasks(includeCompleted bool) []domain.Task
	SaveCycle(domain.Cycle)
	AllCycles() []domain.Cycle
	SaveChanges() error
	SaveSettings(domain.Settings)
	LoadRepository() error
	// Add the following
	// - Zettel
	// - Release
	// - Resource
	// - Artifact
}
