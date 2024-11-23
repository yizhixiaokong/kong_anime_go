package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

// MovieDAO 定义电影DAO
type MovieDAO struct {
	db *gorm.DB
}

// NewMovieDAO 创建电影DAO
func NewMovieDAO(db *gorm.DB) *MovieDAO {
	return &MovieDAO{db: db}
}

// Create 创建一个新的电影
func (dao *MovieDAO) Create(movie *models.Movie) error {
	return dao.db.Create(movie).Error
}

// GetByID 根据ID获取电影
func (dao *MovieDAO) GetByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	err := dao.db.Preload("Categories").Preload("Tags").First(&movie, id).Error
	return &movie, err
}

// GetAll 获取所有电影
func (dao *MovieDAO) GetAll() ([]models.Movie, error) {
	var movies []models.Movie
	err := dao.db.Preload("Categories").Preload("Tags").Find(&movies).Error
	return movies, err
}

// Update 更新电影
func (dao *MovieDAO) Update(movie *models.Movie) error {
	return dao.db.Save(movie).Error
}

// Delete 删除电影
func (dao *MovieDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Movie{}, id).Error
}

// HardDelete 硬删除电影
func (dao *MovieDAO) HardDelete(id uint) error {
	return dao.db.Unscoped().Delete(&models.Movie{}, id).Error
}

// GetByCategory 根据分类获取电影
func (dao *MovieDAO) GetByCategory(categoryName string) ([]models.Movie, error) {
	var movies []models.Movie
	err := dao.db.Joins("JOIN movie_categories ON movie_categories.movie_id = movies.id").
		Joins("JOIN categories ON categories.id = movie_categories.category_id").
		Where("categories.name = ?", categoryName).
		Preload("Categories").Preload("Tags").
		Find(&movies).Error
	return movies, err
}

// GetByTag 根据标签获取电影
func (dao *MovieDAO) GetByTag(tagName string) ([]models.Movie, error) {
	var movies []models.Movie
	err := dao.db.Joins("JOIN movie_tags ON movie_tags.movie_id = movies.id").
		Joins("JOIN tags ON tags.id = movie_tags.tag_id").
		Where("tags.name = ?", tagName).
		Preload("Categories").Preload("Tags").
		Find(&movies).Error
	return movies, err
}
