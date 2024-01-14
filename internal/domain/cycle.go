package domain

import (
	"gorm.io/gorm"
	"time"
)

type Cycle struct {
	gorm.Model
	Title           string
	Active          bool
	StartDate       time.Time
	Accomplishments []Accomplishment
}

func (c *Cycle) FilterValue() string {
	return c.Title
}
