package models

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Name        string     `gorm:"not null"`                    // 名称
	Categories  []Category `gorm:"many2many:movie_categories;"` // 分类
	Tags        []Tag      `gorm:"many2many:movie_tags;"`       // 标签
	ReleaseYear int        // 上映年份
}
