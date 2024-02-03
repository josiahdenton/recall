package domain

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Artifact struct {
	gorm.Model
	Name string
	Tags string
	// Path is the relative path to the working dir for this Artifact
	Path string
	// Editor is the set editor (run from the cmd line) to open this project in
	Editor    string
	Releases  []Release
	Resources []Resource
}

func (a *Artifact) FilterValue() string {
	return a.Name + a.Tags
}

func (a *Artifact) Open() *exec.Cmd {
	path, err := a.expandPath()
	if err != nil {
		log.Printf("failed to open artifact during path expansion: %v", err)
		return nil
	} else if len(path) < 1 && len(a.Editor) < 1 {
		return nil
	}
	cmd := exec.Command(a.Editor, path)
	cmd.Dir = path
	return cmd
}

func (a *Artifact) expandPath() (string, error) {
	path := a.Path
	if strings.Contains(a.Path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to expand home dir %v", home)
		}
		path = strings.Replace(a.Path, "~", home, 1)
	}
	return path, nil
}
