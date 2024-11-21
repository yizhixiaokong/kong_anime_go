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

func (dao *AnimeDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Anime{}, id).Error
}

func (dao *AnimeDAO) Update(anime *models.Anime) error {
	return dao.db.Save(anime).Error
}

func (dao *AnimeDAO) GetByID(id uint) (*models.Anime, error) {
	var anime models.Anime
	err := dao.db.Preload("Categories").Preload("Tags").First(&anime, id).Error
	return &anime, err
}

func (dao *AnimeDAO) GetAllPaginated(page, pageSize int) ([]models.Anime, int64, error) {
	var animes []models.Anime
	var total int64
	offset := (page - 1) * pageSize
	err := dao.db.Preload("Categories").Preload("Tags").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Count(&total)
	return animes, total, err
}

func (dao *AnimeDAO) GetByNameAndAlias(name string, page, pageSize int) ([]models.Anime, int64, error) {
	var animes []models.Anime
	var total int64
	offset := (page - 1) * pageSize
	err := dao.db.Where("name LIKE ? OR aliases LIKE ?", "%"+name+"%", "%"+name+"%").
		Preload("Categories").Preload("Tags").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Where("name LIKE ? OR aliases LIKE ?", "%"+name+"%", "%"+name+"%").Count(&total)
	return animes, total, err
}

func (dao *AnimeDAO) GetBySeason(season string, page, pageSize int) ([]models.Anime, int64, error) {
	var animes []models.Anime
	var total int64
	offset := (page - 1) * pageSize
	err := dao.db.Where("season = ?", season).
		Preload("Categories").Preload("Tags").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Where("season = ?", season).Count(&total)
	return animes, total, err
}

func (dao *AnimeDAO) GetByCategory(categoryName string, page, pageSize int) ([]models.Anime, int64, error) {
	var animes []models.Anime
	var total int64
	offset := (page - 1) * pageSize
	err := dao.db.Joins("JOIN anime_categories ON anime_categories.anime_id = animes.id").
		Joins("JOIN categories ON categories.id = anime_categories.category_id").
		Where("categories.name = ?", categoryName).
		Preload("Categories").Preload("Tags").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Joins("JOIN anime_categories ON anime_categories.anime_id = animes.id").
		Joins("JOIN categories ON categories.id = anime_categories.category_id").
		Where("categories.name = ?", categoryName).Count(&total)
	return animes, total, err
}

func (dao *AnimeDAO) GetByTag(tagName string, page, pageSize int) ([]models.Anime, int64, error) {
	var animes []models.Anime
	var total int64
	offset := (page - 1) * pageSize
	err := dao.db.Joins("JOIN anime_tags ON anime_tags.anime_id = animes.id").
		Joins("JOIN tags ON tags.id = anime_tags.tag_id").
		Where("tags.name = ?", tagName).
		Preload("Categories").Preload("Tags").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Joins("JOIN anime_tags ON anime_tags.anime_id = animes.id").
		Joins("JOIN tags ON tags.id = anime_tags.tag_id").
		Where("tags.name = ?", tagName).Count(&total)
	return animes, total, err
}
