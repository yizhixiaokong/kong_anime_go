package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name   string  `gorm:"unique;not null"`
	Animes []Anime `gorm:"many2many:anime_tags;"`
	Movies []Movie `gorm:"many2many:movie_tags;"`
}
