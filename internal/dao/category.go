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

func (dao *CategoryDAO) GetByName(name string) (*models.Category, error) {
	var category models.Category
	err := dao.db.Where("name = ?", name).First(&category).Error
	return &category, err
}

func (dao *CategoryDAO) GetByNameLike(name string) ([]models.Category, error) {
	var categories []models.Category
	err := dao.db.Where("name LIKE ?", "%"+name+"%").Find(&categories).Error
	return categories, err
}

func (dao *CategoryDAO) CheckRelatedItems(id uint) (bool, error) {
	var animeCount, movieCount int64
	err := dao.db.Model(&models.Anime{}).Where("id IN (SELECT anime_id FROM anime_categories WHERE category_id = ?)", id).Count(&animeCount).Error
	if err != nil {
		return false, err
	}
	err = dao.db.Model(&models.Movie{}).Where("id IN (SELECT movie_id FROM movie_categories WHERE category_id = ?)", id).Count(&movieCount).Error
	if err != nil {
		return false, err
	}
	return animeCount > 0 || movieCount > 0, nil
}
