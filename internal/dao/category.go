package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

// CategoryDAO 定义分类DAO
type CategoryDAO struct {
	db *gorm.DB
}

// NewCategoryDAO 创建分类DAO
func NewCategoryDAO(db *gorm.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

// Create 创建一个新的分类
func (dao *CategoryDAO) Create(category *models.Category) error {
	return dao.db.Create(category).Error
}

// GetByID 根据ID获��分类
func (dao *CategoryDAO) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := dao.db.First(&category, id).Error
	return &category, err
}

// GetAll 获取所有分类
func (dao *CategoryDAO) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := dao.db.Find(&categories).Error
	return categories, err
}

// Update 更新分类
func (dao *CategoryDAO) Update(category *models.Category) error {
	return dao.db.Save(category).Error
}

// Delete 删除分类
func (dao *CategoryDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Category{}, id).Error
}

// HardDelete 硬删除分类
func (dao *CategoryDAO) HardDelete(id uint) error {
	return dao.db.Unscoped().Delete(&models.Category{}, id).Error
}

// GetByName 根据名称获取分类
func (dao *CategoryDAO) GetByName(name string) (*models.Category, error) {
	var category models.Category
	err := dao.db.Where("name = ?", name).First(&category).Error
	return &category, err
}

// GetByNameLike 根据名称模糊查询分类
func (dao *CategoryDAO) GetByNameLike(name string) ([]models.Category, error) {
	var categories []models.Category
	err := dao.db.Where("name LIKE ?", "%"+name+"%").Find(&categories).Error
	return categories, err
}

// CheckRelatedItems 检查分类是否关联了其他项目
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

// GetCategoryStats 获取分类统计信息
func (dao *CategoryDAO) GetCategoryStats() (map[string]int, error) {
	var results []struct {
		Name  string
		Count int
	}
	err := dao.db.Table("categories").
		Select("categories.name as name, COUNT(anime_categories.anime_id) as count").
		Joins("JOIN anime_categories ON anime_categories.category_id = categories.id").
		Group("categories.name").
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
