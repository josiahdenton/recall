package domain

import (
	"github.com/google/uuid"
	"os/exec"
	"runtime"
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
	return r.Name
}

func (r *Resource) OpenLink() bool {
	if r.Type != WebResource {
		return false
	}
	url := r.Source
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
