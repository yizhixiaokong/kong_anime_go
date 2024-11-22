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
		Order("id DESC").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Count(&total)
	return animes, total, err
}

func (dao *AnimeDAO) GetByNameAndAlias(name string, page, pageSize int) ([]models.Anime, int64, error) {
	return dao.getByCondition("name LIKE ? OR aliases LIKE ?", "%"+name+"%", page, pageSize)
}

func (dao *AnimeDAO) GetBySeason(season string, page, pageSize int) ([]models.Anime, int64, error) {
	return dao.getByCondition("season = ?", season, page, pageSize)
}

func (dao *AnimeDAO) GetByCategory(categoryName string, page, pageSize int) ([]models.Anime, int64, error) {
	return dao.getByJoinCondition("categories.name = ?", categoryName, "anime_categories", "categories", "category_id", page, pageSize)
}

func (dao *AnimeDAO) GetByTag(tagName string, page, pageSize int) ([]models.Anime, int64, error) {
	return dao.getByJoinCondition("tags.name = ?", tagName, "anime_tags", "tags", "tag_id", page, pageSize)
}

func (dao *AnimeDAO) ClearCategories(animeID uint) error {
	return dao.clearAssociations(animeID, "Categories")
}

func (dao *AnimeDAO) ClearTags(animeID uint) error {
	return dao.clearAssociations(animeID, "Tags")
}

func (dao *AnimeDAO) AddCategory(animeID, categoryID uint) error {
	a := models.Anime{}
	a.ID = animeID
	category := models.Category{}
	category.ID = categoryID
	return dao.db.Model(&a).Association("Categories").Append(&category)
}

func (dao *AnimeDAO) AddTag(animeID, tagID uint) error {
	a := models.Anime{}
	a.ID = animeID
	tag := models.Tag{}
	tag.ID = tagID
	return dao.db.Model(&a).Association("Tags").Append(&tag)
}

type SeasonCount struct {
	Season string
	Count  int
}

func (dao *AnimeDAO) GetAllSeasons() ([]SeasonCount, error) {
	var seasons []SeasonCount
	err := dao.db.Model(&models.Anime{}).
		Select("season, COUNT(*) as count").
		Group("season").
		Order("season").
		Scan(&seasons).Error
	return seasons, err
}

func (dao *AnimeDAO) getByCondition(condition string, args interface{}, page, pageSize int) ([]models.Anime, int64, error) {
	var animes []models.Anime
	var total int64
	offset := (page - 1) * pageSize
	err := dao.db.Where(condition, args).
		Preload("Categories").Preload("Tags").
		Order("id DESC").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Where(condition, args).Count(&total)
	return animes, total, err
}

func (dao *AnimeDAO) getByJoinCondition(condition, value, joinTable, joinModel, joinField string, page, pageSize int) ([]models.Anime, int64, error) {
	var animes []models.Anime
	var total int64
	offset := (page - 1) * pageSize
	err := dao.db.Joins("JOIN "+joinTable+" ON "+joinTable+".anime_id = animes.id").
		Joins("JOIN "+joinModel+" ON "+joinModel+".id = "+joinTable+"."+joinField).
		Where(condition, value).
		Preload("Categories").Preload("Tags").
		Order("id DESC").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Joins("JOIN "+joinTable+" ON "+joinTable+".anime_id = animes.id").
		Joins("JOIN "+joinModel+" ON "+joinModel+".id = "+joinTable+"."+joinField).
		Where(condition, value).Count(&total)
	return animes, total, err
}

func (dao *AnimeDAO) clearAssociations(animeID uint, association string) error {
	a := models.Anime{}
	a.ID = animeID
	return dao.db.Model(&a).Association(association).Clear()
}
