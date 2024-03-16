package repository

import (
	"github.com/josiahdenton/recall/internal/domain"
)

type Repository interface {

	// Tasks

	Task(uint) *domain.Task
	DeleteTask(uint)
	ModifyTask(domain.Task)
	AllTasks() []domain.Task
	ArchivedTasks() []domain.Task
	UnlinkTaskResource(*domain.Task, *domain.Resource)
	UnlinkTaskStep(*domain.Task, *domain.Step)
	UnlinkTaskStatus(*domain.Task, *domain.Status)
	UndoDeleteTask(uint)

	// Accomplishments

	Cycle(uint) *domain.Cycle // etc...
	Accomplishment(uint) *domain.Accomplishment
	ModifyAccomplishment(domain.Accomplishment)
	DeleteAccomplishment(uint)
	UndoDeleteAccomplishment(uint)
	ModifyCycle(domain.Cycle)
	AllCycles() []domain.Cycle

	// Resources

	ModifyResource(domain.Resource)
	AllResources() []domain.Resource

	// Status

	ModifyStatus(domain.Status)

	// Step

	ModifyStep(step domain.Step)

	LoadRepository() error
	// Add Delete*FromTask for Resource, Status, Step
	//Resource(string) domain.Resource // I don't think I ever would need to get a single resource
	// Add the following
	// - Zettel
	// - Release
	// - Resource
	// - Artifact
}
