package domain

import (
	"gorm.io/gorm"
)

type Zettel struct {
	gorm.Model
	Name      string
	Tags      string
	Concept   string
	Favorite  bool
	Links     []*Zettel  `gorm:"many2many:zettels_zettels"`
	Resources []Resource `gorm:"many2many:zettel_resources"`
}

func (z *Zettel) FilterValue() string {
	return z.Name + z.Tags
}
