package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name   string  `gorm:"unique;not null"`
	Animes []Anime `gorm:"many2many:anime_categories;" json:"omitempty"`
	Movies []Movie `gorm:"many2many:movie_categories;" json:"omitempty"`
}
