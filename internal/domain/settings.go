package domain

import "gorm.io/gorm"

type Settings struct {
	gorm.Model
	Editor string
}
