package domain

import (
	"gorm.io/gorm"
)

const (
	NoneType = iota
	WebResource
	ZettelResource // TODO when you link a Zettel, you'll need to create some sort of lookup
	FilePathResource
	TaskResource
)

type ResourceType = int

type Resource struct {
	gorm.Model
	Name   string
	Source string
	Tags   string
	Type   ResourceType
	TaskID uint
}

func (r *Resource) StringType() string {
	switch r.Type {
	case WebResource:
		return "Web"
	case ZettelResource:
		return "Zettel"
	case FilePathResource:
		return "File"
	case TaskResource:
		return "Task"
	default:
		return ""
	}
}

func (r *Resource) FilterValue() string {
	return r.Name + r.Tags
}

func (r *Resource) Open() bool {
	if r.Type != WebResource {
		return false
	}
	return openWebPage(r.Source)
}
