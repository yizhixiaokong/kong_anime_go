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

func (dao *TagDAO) GetByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := dao.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

func (dao *TagDAO) GetByNameLike(name string) ([]models.Tag, error) {
	var tags []models.Tag
	err := dao.db.Where("name LIKE ?", "%"+name+"%").Find(&tags).Error
	return tags, err
}

func (dao *TagDAO) CheckRelatedItems(id uint) (bool, error) {
	var animeCount, movieCount int64
	err := dao.db.Model(&models.Anime{}).Where("id IN (SELECT anime_id FROM anime_tags WHERE tag_id = ?)", id).Count(&animeCount).Error
	if err != nil {
		return false, err
	}
	err = dao.db.Model(&models.Movie{}).Where("id IN (SELECT movie_id FROM movie_tags WHERE tag_id = ?)", id).Count(&movieCount).Error
	if err != nil {
		return false, err
	}
	return animeCount > 0 || movieCount > 0, nil
}

func (dao *TagDAO) GetTagStats() (map[string]int, error) {
	var results []struct {
		Name  string
		Count int
	}
	err := dao.db.Table("tags").
		Select("tags.name as name, COUNT(anime_tags.anime_id) as count").
		Joins("JOIN anime_tags ON anime_tags.tag_id = tags.id").
		Group("tags.name").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	for _, result := range results {
		stats[result.Name] = result.Count
	}
	return stats, nil
}
