package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

// AnimeDAO 定义动漫DAO
type AnimeDAO struct {
	db *gorm.DB
}

// NewAnimeDAO 创建动漫DAO
func NewAnimeDAO(db *gorm.DB) *AnimeDAO {
	return &AnimeDAO{db: db}
}

// Create 创建一个新的动漫
func (dao *AnimeDAO) Create(anime *models.Anime) error {
	return dao.db.Create(anime).Error
}

// Delete 删除动漫
func (dao *AnimeDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Anime{}, id).Error
}

// HardDelete 硬删除动漫
func (dao *AnimeDAO) HardDelete(id uint) error {
	return dao.db.Unscoped().Delete(&models.Anime{}, id).Error
}

// Update 更新动漫
func (dao *AnimeDAO) Update(anime *models.Anime) error {
	return dao.db.Save(anime).Error
}

// GetByID 根据ID获取动漫
func (dao *AnimeDAO) GetByID(id uint) (*models.Anime, error) {
	var anime models.Anime
	err := dao.db.Preload("Categories").Preload("Tags").First(&anime, id).Error
	return &anime, err
}

// GetAllPaginated 获取分页的动漫列表
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

// GetByNameAndAlias 根据名称和别名获取动漫
func (dao *AnimeDAO) GetByNameAndAlias(name string, page, pageSize int) ([]models.Anime, int64, error) {
	condition := "name LIKE ? OR aliases LIKE ?"
	args := []any{"%" + name + "%", "%" + name + "%"}
	return dao.getByCondition(condition, args, page, pageSize)
}

// GetBySeason 根据季节获取动漫
func (dao *AnimeDAO) GetBySeason(season string, page, pageSize int) ([]models.Anime, int64, error) {
	return dao.getByCondition("season = ?", []any{season}, page, pageSize)
}

// GetByCategory 根据分类获取动漫
func (dao *AnimeDAO) GetByCategory(categoryName string, page, pageSize int) ([]models.Anime, int64, error) {
	return dao.getByJoinCondition("categories.name = ?", categoryName, "anime_categories", "categories", "category_id", page, pageSize)
}

// GetByTag 根据标签获取动漫
func (dao *AnimeDAO) GetByTag(tagName string, page, pageSize int) ([]models.Anime, int64, error) {
	return dao.getByJoinCondition("tags.name = ?", tagName, "anime_tags", "tags", "tag_id", page, pageSize)
}

// ClearCategories 清除动漫的所有分类
func (dao *AnimeDAO) ClearCategories(animeID uint) error {
	return dao.clearAssociations(animeID, "Categories")
}

// ClearTags 清除动漫的所有标签
func (dao *AnimeDAO) ClearTags(animeID uint) error {
	return dao.clearAssociations(animeID, "Tags")
}

// AddCategory 为动漫添加分类
func (dao *AnimeDAO) AddCategory(animeID, categoryID uint) error {
	a := models.Anime{}
	a.ID = animeID
	category := models.Category{}
	category.ID = categoryID
	return dao.db.Model(&a).Association("Categories").Append(&category)
}

// AddTag 为动漫添加标签
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

// GetAllSeasons 获取所有季节及其动漫数量
func (dao *AnimeDAO) GetAllSeasons() ([]SeasonCount, error) {
	var seasons []SeasonCount
	err := dao.db.Model(&models.Anime{}).
		Select("season, COUNT(*) as count").
		Group("season").
		Order("season").
		Scan(&seasons).Error
	return seasons, err
}

func (dao *AnimeDAO) GetFollowsByAnimeID(animeID uint) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64
	err := dao.db.Where("anime_id = ?", animeID).Find(&follows).Count(&total).Error
	return follows, total, err
}

func (dao *AnimeDAO) getByCondition(condition string, args []any, page, pageSize int) ([]models.Anime, int64, error) {
	var animes []models.Anime
	var total int64
	offset := (page - 1) * pageSize
	err := dao.db.Where(condition, args...).
		Preload("Categories").Preload("Tags").
		Order("id DESC").
		Limit(pageSize).Offset(offset).
		Find(&animes).Error
	dao.db.Model(&models.Anime{}).Where(condition, args...).Count(&total)
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
