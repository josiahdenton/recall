package domain

import "gorm.io/gorm"

type Artifact struct {
	gorm.Model
	Name      string
	Tags      string
	Releases  []Release
	Resources []Resource
}

func (a *Artifact) FilterValue() string {
	return a.Name
}
