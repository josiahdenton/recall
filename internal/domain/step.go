package domain

import (
	"gorm.io/gorm"
)

type Step struct {
	gorm.Model
	Description string
	Complete    bool
	TaskID      uint
	ReleaseID   uint
}

func (s *Step) ToggleStatus() {
	s.Complete = !s.Complete
}

func (s *Step) FilterValue() string {
	return s.Description
}
