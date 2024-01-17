package domain

import (
	"gorm.io/gorm"
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
	gorm.Model
	Name     string
	Source   string
	Type     ResourceType
	TaskID   uint
	ZettelID uint
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

func (r *Resource) Open() bool {
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
