package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

type MovieDAO struct {
	db *gorm.DB
}

func NewMovieDAO(db *gorm.DB) *MovieDAO {
	return &MovieDAO{db: db}
}

func (dao *MovieDAO) Create(movie *models.Movie) error {
	return dao.db.Create(movie).Error
}

func (dao *MovieDAO) GetByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	err := dao.db.Preload("Categories").Preload("Tags").First(&movie, id).Error
	return &movie, err
}

func (dao *MovieDAO) GetAll() ([]models.Movie, error) {
	var movies []models.Movie
	err := dao.db.Preload("Categories").Preload("Tags").Find(&movies).Error
	return movies, err
}

func (dao *MovieDAO) Update(movie *models.Movie) error {
	return dao.db.Save(movie).Error
}

func (dao *MovieDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Movie{}, id).Error
}

func (dao *MovieDAO) GetByCategory(categoryName string) ([]models.Movie, error) {
	var movies []models.Movie
	err := dao.db.Joins("JOIN movie_categories ON movie_categories.movie_id = movies.id").
		Joins("JOIN categories ON categories.id = movie_categories.category_id").
		Where("categories.name = ?", categoryName).
		Preload("Categories").Preload("Tags").
		Find(&movies).Error
	return movies, err
}

func (dao *MovieDAO) GetByTag(tagName string) ([]models.Movie, error) {
	var movies []models.Movie
	err := dao.db.Joins("JOIN movie_tags ON movie_tags.movie_id = movies.id").
		Joins("JOIN tags ON tags.id = movie_tags.tag_id").
		Where("tags.name = ?", tagName).
		Preload("Categories").Preload("Tags").
		Find(&movies).Error
	return movies, err
}
