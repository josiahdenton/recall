package repository

import (
	"github.com/josiahdenton/recall/internal/domain"
)

type Repository interface {
	Task(uint) *domain.Task
	DeleteTask(uint)
	Cycle(uint) *domain.Cycle // etc...
	Accomplishment(uint) *domain.Accomplishment
	ModifyAccomplishment(domain.Accomplishment)
	ModifyTask(domain.Task)
	AllTasks() []domain.Task
	ArchivedTasks() []domain.Task
	ModifyCycle(domain.Cycle)
	AllCycles() []domain.Cycle
	ModifyResource(domain.Resource)
	AllResources() []domain.Resource
	ModifySettings(domain.Settings)
	AllZettels() []domain.Zettel
	Zettel(uint) *domain.Zettel
	ModifyZettel(domain.Zettel)
	LoadRepository() error
	//Resource(string) domain.Resource // I don't think I ever would need to get a single resource
	// Add the following
	// - Zettel
	// - Release
	// - Resource
	// - Artifact
}
