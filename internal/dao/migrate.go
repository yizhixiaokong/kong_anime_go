package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

// Migrate 迁移数据库
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Anime{}, &models.Category{}, &models.Tag{}, &models.Movie{})
}
