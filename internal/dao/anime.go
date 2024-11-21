package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

type AnimeDAO struct {
	db *gorm.DB
}

func NewAnimeDAO(db *gorm.DB) *AnimeDAO {
	return &AnimeDAO{db: db}
}

func (dao *AnimeDAO) Create(anime *models.Anime) error {
	return dao.db.Create(anime).Error
}

func (dao *AnimeDAO) GetByID(id uint) (*models.Anime, error) {
	var anime models.Anime
	err := dao.db.Preload("Categories").Preload("Tags").First(&anime, id).Error
	return &anime, err
}

func (dao *AnimeDAO) GetAll() ([]models.Anime, error) {
	var animes []models.Anime
	err := dao.db.Preload("Categories").Preload("Tags").Find(&animes).Error
	return animes, err
}

func (dao *AnimeDAO) Update(anime *models.Anime) error {
	return dao.db.Save(anime).Error
}

func (dao *AnimeDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Anime{}, id).Error
}

func (dao *AnimeDAO) GetByAlias(alias string) ([]models.Anime, error) {
	var animes []models.Anime
	err := dao.db.Where("aliases LIKE ?", "%"+alias+"%").
		Preload("Categories").Preload("Tags").
		Find(&animes).Error
	return animes, err
}

func (dao *AnimeDAO) GetByCategory(categoryName string) ([]models.Anime, error) {
	var animes []models.Anime
	err := dao.db.Joins("JOIN anime_categories ON anime_categories.anime_id = animes.id").
		Joins("JOIN categories ON categories.id = anime_categories.category_id").
		Where("categories.name = ?", categoryName).
		Preload("Categories").Preload("Tags").
		Find(&animes).Error
	return animes, err
}

func (dao *AnimeDAO) GetByTag(tagName string) ([]models.Anime, error) {
	var animes []models.Anime
	err := dao.db.Joins("JOIN anime_tags ON anime_tags.anime_id = animes.id").
		Joins("JOIN tags ON tags.id = anime_tags.tag_id").
		Where("tags.name = ?", tagName).
		Preload("Categories").Preload("Tags").
		Find(&animes).Error
	return animes, err
}
