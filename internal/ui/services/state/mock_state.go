package state

import (
	"github.com/josiahdenton/recall/internal/domain"
	"time"
)

var (
	mockTask = domain.Task{
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
	}
	mockTasks = []domain.Task{
		{
			Title: "Take out Garbage",
			Tags:  "13",
		},
		{
			Title: "Take out your mom",
			Tags:  "03",
		},
		{
			Title: "finish the task app",
			Tags:  "103",
		},
	}
	mockResources = []domain.Resource{
		{
			Name:   "Google",
			Source: "https://www.google.com",
			Tags:   "t1,t2,t3,t4",
			Type:   domain.WebResource,
		},
		{
			Name:   "Youtube",
			Source: "https://www.youtube.com/",
			Tags:   "t1,t2,t3,t4",
			Type:   domain.WebResource,
		},
	}
)
