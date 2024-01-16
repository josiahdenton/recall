package domain

import (
	"gorm.io/gorm"
	"strings"
)

type Zettel struct {
	gorm.Model
	Name            string
	ContentLocation string    // an absolute path to a file
	Links           []*Zettel `gorm:"many2many:zettels_zettels"` // TODO this may not be right
}

func (z *Zettel) FilterValue() string {
	return z.Name
}

func NewZettel(name string) Zettel {
	name = strings.ToLower(name)
	filePath := strings.Replace(name, " ", "_", -1)
	filePath += ".md"
	return Zettel{
		Name:            name,
		ContentLocation: filePath,
	}
}
