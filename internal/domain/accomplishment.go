package domain

import "github.com/google/uuid"

type Accomplishment struct {
	Id                string   `json:"id"`
	Description       string   `json:"description"`
	Impact            string   `json:"impact"`
	Strength          string   `json:"strength"`
	AssociatedTaskIds []string `json:"associated_task_ids"`
	associatedTasks   []Task
}

func (a *Accomplishment) FilterValue() string {
	return a.Description
}

func (a *Accomplishment) AssociatedTasks() []Task {
	return a.associatedTasks
}

func (a *Accomplishment) AttachAssociatedTasks(tasks []Task) {
	a.associatedTasks = tasks
}

func NewAccomplishment(description, impact, strength string) Accomplishment {
	id, err := uuid.NewRandom()
	if err != nil {
		return Accomplishment{}
	}

	return Accomplishment{
		Id:                id.String(),
		Description:       description,
		Impact:            impact,
		Strength:          strength,
		AssociatedTaskIds: make([]string, 0),
	}
}
