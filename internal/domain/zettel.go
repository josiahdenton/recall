package domain

import (
	"gorm.io/gorm"
	"strings"
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

// Summary returns the first 20 lines of a zettel
func (z *Zettel) Summary() string {
	var b strings.Builder
	for i, line := range strings.Split(z.Concept, "\n") {
		if i < 20 {
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	return b.String()
}
