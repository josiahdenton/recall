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
	DeleteAccomplishment(uint)
	ModifyTask(domain.Task)
	AllTasks() []domain.Task
	ArchivedTasks() []domain.Task
	ModifyCycle(domain.Cycle)
	AllCycles() []domain.Cycle
	ModifyResource(domain.Resource)
	ModifyStatus(domain.Status)
	AllResources() []domain.Resource
	ModifySettings(domain.Settings)
	AllZettels() []domain.Zettel
	Zettel(uint) *domain.Zettel
	ModifyZettel(domain.Zettel)
	DeleteZettel(uint)
	UnlinkZettel(*domain.Zettel, *domain.Zettel)
	UnlinkZettelResource(*domain.Zettel, *domain.Resource)
	LoadRepository() error
	UnlinkTaskResource(*domain.Task, *domain.Resource)
	UnlinkTaskStep(*domain.Task, *domain.Step)
	UnlinkTaskStatus(*domain.Task, *domain.Status)
	ModifyStep(step domain.Step)
	UndoDeleteTask(uint)
	UndoDeleteZettel(uint)
	UndoDeleteAccomplishment(uint)
	AllArtifacts() []domain.Artifact
	Artifact(uint) *domain.Artifact
	ModifyArtifact(domain.Artifact)
	DeleteArtifact(uint)
	Release(uint) *domain.Release
	ModifyRelease(domain.Release)
	DeleteRelease(uint)
	DeleteArtifactRelease(*domain.Artifact, *domain.Release)
	DeleteArtifactResource(*domain.Artifact, *domain.Resource)
	UndoDeleteArtifact(uint)
	// Add Delete*FromTask for Resource, Status, Step
	//Resource(string) domain.Resource // I don't think I ever would need to get a single resource
	// Add the following
	// - Zettel
	// - Release
	// - Resource
	// - Artifact
}
