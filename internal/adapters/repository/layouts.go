package repository

import "github.com/josiahdenton/recall/internal/domain"

type cyclesLayout struct {
	Cycles []domain.Cycle `json:"cycles"`
}

type accomplishmentsLayout struct {
	Accomplishments map[string]domain.Accomplishment `json:"accomplishments"`
}

type tasksLayout struct {
	Tasks map[string]domain.Task `json:"tasks"`
}

type settingsLayout struct {
	Settings domain.Settings
}
