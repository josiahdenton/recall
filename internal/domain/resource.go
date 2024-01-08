package domain

import "github.com/google/uuid"

const (
	WebResource    = iota
	ZettelResource // TODO when you link a Zettel, you'll need to create some sort of lookup
	FilePathResource
	TaskResource
)

type ResourceType = int

type Resource struct {
	Id     string       `json:"id"`
	Name   string       `json:"name"`
	Source string       `json:"source"`
	Type   ResourceType `json:"type"`
}

func NewResource(name, source string, resourceType ResourceType) Resource {
	id, err := uuid.NewRandom()
	if err != nil {
		return Resource{}
	}

	return Resource{
		Id:     id.String(),
		Name:   name,
		Source: source,
		Type:   resourceType,
	}
}

func (r *Resource) FilterValue() string {
	return r.Name
}
