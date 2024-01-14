package domain

import (
	"gorm.io/gorm"
)

type Accomplishment struct {
	gorm.Model
	Description string
	Impact      string
	Strength    string
	Tasks       []Task
	CycleID     uint
}

func (a *Accomplishment) FilterValue() string {
	return a.Description
}
