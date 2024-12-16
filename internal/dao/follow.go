package dao

import (
	"kong-anime-go/internal/dao/models"

	"gorm.io/gorm"
)

// FollowDAO 定义追番DAO
type FollowDAO struct {
	db *gorm.DB
}

// NewFollowDAO 创建追番DAO
func NewFollowDAO(db *gorm.DB) *FollowDAO {
	return &FollowDAO{db: db}
}

// Create 创建一个新的追番
func (dao *FollowDAO) Create(follow *models.Follow) error {
	return dao.db.Create(follow).Error
}

// Delete 删除追番
func (dao *FollowDAO) Delete(id uint) error {
	return dao.db.Delete(&models.Follow{}, id).Error
}

// HardDelete 硬删除追番
func (dao *FollowDAO) HardDelete(id uint) error {
	return dao.db.Unscoped().Delete(&models.Follow{}, id).Error
}

// Update 更新追番
func (dao *FollowDAO) Update(follow *models.Follow) error {
	return dao.db.Save(follow).Error
}

// GetByID 根据ID获取追番
func (dao *FollowDAO) GetByID(id uint) (*models.Follow, error) {
	var follow models.Follow
	err := dao.db.Preload("Anime").First(&follow, id).Error
	return &follow, err
}

// GetByAnimeID 根据AnimeID获取追番
func (dao *FollowDAO) GetByAnimeID(animeID uint) (*models.Follow, error) {
	var follow models.Follow
	err := dao.db.Where("anime_id = ?", animeID).First(&follow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &follow, err
}

// GetAllPaginated 获取分页的追番列表
func (dao *FollowDAO) GetAllPaginated(page, pageSize int, category int, status int, name *string, sorter *string) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64
	offset := (page - 1) * pageSize
	query := dao.db.Preload("Anime").Limit(pageSize).Offset(offset)
	countQuery := dao.db.Model(&models.Follow{})
	if category != -1 {
		query = query.Where("category = ?", category)
		countQuery = countQuery.Where("category = ?", category)
	}
	if status != -1 {
		query = query.Where("status = ?", status)
		countQuery = countQuery.Where("status = ?", status)
	}
	query = query.Joins("JOIN animes ON animes.id = follows.anime_id")
	countQuery = countQuery.Joins("JOIN animes ON animes.id = follows.anime_id")

	if name != nil && *name != "" {
		query = query.Where("animes.name LIKE ?", "%"+*name+"%")
		countQuery = countQuery.Where("animes.name LIKE ?", "%"+*name+"%")
	}
	if sorter != nil && *sorter != "" {
		query = query.Order(*sorter)
	} else {
		query = query.Order("id DESC")
	}
	err := query.Find(&follows).Error
	countQuery.Count(&total)
	return follows, total, err
}

// GetFollowsByAnimeID 根据AnimeID获取追番
func (dao *FollowDAO) GetFollowsByAnimeID(animeID uint) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64
	err := dao.db.Where("anime_id = ?", animeID).Find(&follows).Count(&total).Error
	return follows, total, err
}

// HardDeleteByAnimeID
func (dao *FollowDAO) HardDeleteByAnimeID(animeID uint) error {
	return dao.db.Unscoped().Where("anime_id = ?", animeID).Delete(&models.Follow{}).Error
}
