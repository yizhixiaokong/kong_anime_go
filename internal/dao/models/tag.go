package models

import (
	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	gorm.Model
	Name   string  `gorm:"unique;not null"`
	Animes []Anime `gorm:"many2many:anime_tags;" json:"omitempty"`
	Movies []Movie `gorm:"many2many:movie_tags;" json:"omitempty"`
}
