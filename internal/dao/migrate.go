package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Anime{}, &models.Category{}, &models.Tag{}, &models.Movie{})
}
