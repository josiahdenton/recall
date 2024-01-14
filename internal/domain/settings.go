package domain

import "gorm.io/gorm"

const (
	AccomplishmentsFileName = "accomplishments.json"
	TasksFileName           = "tasks.json"
	CyclesFileName          = "cycles.json"
	SettingsFileName        = "settings.json"
	ResourcesFileName       = "resources.json"
)

type Settings struct {
	gorm.Model
	Location string
}
