package domain

import (
	"gorm.io/gorm"
)

type Zettel struct {
	gorm.Model
	Name      string
	Concept   string
	Favorite  bool
	Links     []*Zettel  `gorm:"many2many:zettels_zettels"` // TODO this may not be right
	Resources []Resource `gorm:"many2many:zettel_resources"`
}

func (z *Zettel) FilterValue() string {
	return z.Name
}
