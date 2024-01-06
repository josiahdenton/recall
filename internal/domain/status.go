package domain

import "github.com/google/uuid"

type Status struct {
	Id          string
	Description string
}

func NewStatus(description string) Status {
	id, err := uuid.NewRandom()
	if err != nil {
		return Status{}
	}

	return Status{
		Id:          id.String(),
		Description: description,
	}
}

func (s *Status) FilterValue() string {
	return s.Description
}
