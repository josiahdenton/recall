package domain

import "github.com/google/uuid"

type Accomplishment struct {
	Id              string
	Description     string
	Impact          string
	Strength        string
	AssociatedTasks []Task // any child / nested data like this is filled in by the repository as needed...
}

func (a *Accomplishment) FilterValue() string {
	return a.Description
}

func NewAccomplishment(description, impact, strength string) Accomplishment {
	id, err := uuid.NewRandom()
	if err != nil {
		return Accomplishment{}
	}

	return Accomplishment{
		Id:              id.String(),
		Description:     description,
		Impact:          impact,
		Strength:        strength,
		AssociatedTasks: make([]Task, 0),
	}
}
