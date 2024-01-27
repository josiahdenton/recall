package domain

import (
	"gorm.io/gorm"
)

type Status struct {
	gorm.Model
	Description string
	TaskID      uint
}

func (s *Status) FilterValue() string {
	return s.Description
}
