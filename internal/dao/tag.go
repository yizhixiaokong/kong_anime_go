package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

type TagDAO struct {
	db *gorm.DB
}

func NewTagDAO(db *gorm.DB) *TagDAO {
	return &TagDAO{db: db}
}

func (dao *TagDAO) Create(tag *models.Tag) error {
	return dao.db.Create(tag).Error
}

func (dao *TagDAO) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := dao.db.First(&tag, id).Error
	return &tag, err
}

func (dao *TagDAO) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := dao.db.Find(&tags).Error
	return tags, err
}

func (dao *TagDAO) Update(tag *models.Tag) error {
	return dao.db.Save(tag).Error
}

func (dao *TagDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Tag{}, id).Error
}
