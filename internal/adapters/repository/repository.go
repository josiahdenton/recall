package repository

import (
	"github.com/josiahdenton/recall/internal/domain"
)

type Repository interface {
	Task(string) *domain.Task
	DeleteTask(string)
	Cycle(string) *domain.Cycle // etc...
	Accomplishment(string) *domain.Accomplishment
	SaveAccomplishment(domain.Accomplishment)
	LinkedAccomplishments([]string) []domain.Accomplishment
	LinkedTasks([]string) []domain.Task
	SaveTask(domain.Task)
	AllTasks() []domain.Task
	ArchivedTasks() []domain.Task
	SaveCycle(domain.Cycle)
	AllCycles() []domain.Cycle
	SaveResource(domain.Resource)
	AllResources() []domain.Resource
	LinkedResources([]string) []domain.Resource
	SaveChanges() error
	SaveSettings(domain.Settings)
	LoadRepository() error
	//ArchiveTask(string) // this can just be on the task...
	//Resource(string) domain.Resource // I don't think I ever would need to get a single resource
	// Add the following
	// - Zettel
	// - Release
	// - Resource
	// - Artifact
}
