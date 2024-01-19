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
	DeleteZettel(uint)
	UnlinkZettel(*domain.Zettel, *domain.Zettel)
	LoadRepository() error
	DeleteTaskResource(*domain.Task, *domain.Resource)
	DeleteTaskStep(*domain.Task, *domain.Step)
	DeleteTaskStatus(*domain.Task, *domain.Status)
	ModifyStep(step domain.Step)
	// Add Delete*FromTask for Resource, Status, Step
	//Resource(string) domain.Resource // I don't think I ever would need to get a single resource
	// Add the following
	// - Zettel
	// - Release
	// - Resource
	// - Artifact
}
