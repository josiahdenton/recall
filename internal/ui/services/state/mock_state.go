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
		Zettels: []domain.Zettel{
			{
				Name:      "Something something abcs",
				Tags:      "here tag, tags cool, tags neat",
				Concept:   "ABC's are boring.",
				Favorite:  false,
				TaskID:    0,
				Links:     nil,
				Resources: nil,
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
)
