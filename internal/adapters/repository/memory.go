package repository

import (
	"github.com/josiahdenton/recall/internal/domain"
	"time"
)

func DefaultInMemory() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: []domain.Task{
			{
				Title:             "Something really cool",
				Tags:              "here are some tags, another one",
				Description:       "This is a task that is necessary to complete. I would put some context here to better understand the purpose and what I am trying to do.",
				Due:               time.Time{},
				Active:            true,
				Duration:          0,
				LastActivatedTime: time.Time{},
				Archive:           false,
				Status: []domain.Status{
					{
						Description: "This is a long status as an example to show on some task as a way to make sure we know what we're doing here.",
					},
					{
						Description: "Something slightly shorter to offer a differing point of view",
					},
				},
				Resources: []domain.Resource{
					{
						Name:   "HahaHa - funny joke",
						Source: "https://www.google.com",
						Tags:   "dad,jokes",
					},
				},
				Steps: []domain.Step{
					{
						Description: "Make Joke",
					},
					{
						Description: "Laugh loudly",
					},
				},
				AccomplishmentID: 0,
			},
			{
				Title:             "Another Thing",
				Tags:              "t1,t2,t3...etc",
				Description:       "This is a task that is necessary to complete. I would put some context here to better understand the purpose and what I am trying to do.",
				Due:               time.Time{},
				Active:            true,
				Duration:          0,
				LastActivatedTime: time.Time{},
				Archive:           false,
				Status: []domain.Status{
					{
						Description: "This is a long status as an example to show on some task as a way to make sure we know what we're doing here.",
					},
					{
						Description: "Something slightly shorter to offer a differing point of view",
					},
				},
				Resources: []domain.Resource{
					{
						Name:   "HahaHa - funny joke",
						Source: "https://www.google.com",
						Tags:   "dad,jokes",
					},
				},
				Steps: []domain.Step{
					{
						Description: "Hahhahahahha",
					},
					{
						Description: "STAY AWAKE",
					},
				},
				AccomplishmentID: 0,
			},
		},
	}
}

type InMemoryStorage struct {
	tasks []domain.Task
}

func (i InMemoryStorage) Task(u uint) *domain.Task {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) AllTasks() []domain.Task {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) ArchivedTasks() []domain.Task {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) ModifyTask(task domain.Task) domain.Task {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) DeleteTask(u uint) {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) UnlinkTaskResource(task *domain.Task, resource *domain.Resource) {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) UnlinkTaskStep(task *domain.Task, step *domain.Step) {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) UnlinkTaskStatus(task *domain.Task, status *domain.Status) {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) UndoDeleteTask(u uint) {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) ModifyStep(step domain.Step) domain.Step {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) Cycle(u uint) *domain.Cycle {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) AllCycles() []domain.Cycle {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) ModifyCycle(cycle domain.Cycle) domain.Cycle {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) Accomplishment(u uint) *domain.Accomplishment {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) ModifyAccomplishment(accomplishment domain.Accomplishment) domain.Accomplishment {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) DeleteAccomplishment(u uint) {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) UndoDeleteAccomplishment(u uint) {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) ModifyResource(resource domain.Resource) domain.Resource {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) AllResources() []domain.Resource {
	//TODO implement me
	panic("implement me")
}

func (i InMemoryStorage) LoadRepository() error {
	//TODO implement me
	panic("implement me")
}
