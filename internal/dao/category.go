package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

type CategoryDAO struct {
	db *gorm.DB
}

func NewCategoryDAO(db *gorm.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

func (dao *CategoryDAO) Create(category *models.Category) error {
	return dao.db.Create(category).Error
}

func (dao *CategoryDAO) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := dao.db.First(&category, id).Error
	return &category, err
}

func (dao *CategoryDAO) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := dao.db.Find(&categories).Error
	return categories, err
}

func (dao *CategoryDAO) Update(category *models.Category) error {
	return dao.db.Save(category).Error
}

func (dao *CategoryDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Category{}, id).Error
}
