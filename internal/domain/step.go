package domain

import "github.com/google/uuid"

type Step struct {
	Id          string
	Description string
	Complete    bool
}

func NewStep(description string) Step {
	id, err := uuid.NewRandom()
	if err != nil {
		return Step{}
	}

	return Step{
		Id:          id.String(),
		Description: description,
	}
}

func (s *Step) ToggleStatus() {
	s.Complete = !s.Complete
}

func (s *Step) FilterValue() string {
	return s.Description
}
