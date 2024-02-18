package state

import (
	"github.com/josiahdenton/recall/internal/domain"
)

var (
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
