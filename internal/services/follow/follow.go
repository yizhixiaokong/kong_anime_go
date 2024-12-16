package follow

import (
	"errors"
	"kong-anime-go/internal/dao"
	"kong-anime-go/internal/dao/models"
)

// Service 处理追番相关的服务
type Service struct {
	followDAO *dao.FollowDAO
	animeDAO  *dao.AnimeDAO
}

// NewService 创建一个新的 FollowService
func NewService(followDAO *dao.FollowDAO, animeDAO *dao.AnimeDAO) *Service {
	return &Service{
		followDAO: followDAO,
		animeDAO:  animeDAO,
	}
}

// Create 创建一个新的追番
func (s *Service) Create(follow *models.Follow) (*models.Follow, error) {
	existingFollow, err := s.followDAO.GetByAnimeID(follow.AnimeID)
	if err != nil {
		return nil, err
	}
	if existingFollow != nil {
		return nil, errors.New("follow already exists for this anime")
	}
	anime, err := s.animeDAO.GetByID(follow.AnimeID)
	if err != nil {
		return nil, err
	}
	if anime == nil {
		return nil, errors.New("anime not found")
	}
	if err := s.followDAO.Create(follow); err != nil {
		return nil, err
	}
	return s.followDAO.GetByID(follow.ID)
}

// Delete 删除一个追番
func (s *Service) Delete(id uint) (uint, error) {
	return id, s.followDAO.Delete(id)
}

// Update 更新一个追番
func (s *Service) Update(follow *models.Follow) (*models.Follow, error) {
	existingFollow, err := s.followDAO.GetByID(follow.ID)
	if err != nil {
		return nil, err
	}
	if existingFollow == nil {
		return nil, errors.New("follow not found")
	}
	if existingFollow.AnimeID != follow.AnimeID {
		return nil, errors.New("cannot change AnimeID")
	}
	existingFollow.Category = follow.Category
	existingFollow.Status = follow.Status
	existingFollow.FinishedAt = follow.FinishedAt
	if err := s.followDAO.Update(existingFollow); err != nil {
		return nil, err
	}
	return s.followDAO.GetByID(follow.ID)
}

// GetByID 根据ID获取追番
func (s *Service) GetByID(id uint) (*models.Follow, error) {
	return s.followDAO.GetByID(id)
}

// GetByAnimeID 根据AnimeID获取追番
func (s *Service) GetByAnimeID(animeID uint) (*models.Follow, error) {
	return s.followDAO.GetByAnimeID(animeID)
}

// GetAll 获取所有追番
func (s *Service) GetAll(page, pageSize int, category *int, status *int, name *string, sorter *string) ([]models.Follow, int64, error) {
	var categoryInt, statusInt int
	if category != nil {
		categoryInt = *category
	} else {
		categoryInt = -1
	}
	if status != nil {
		statusInt = *status
	} else {
		statusInt = -1
	}
	return s.followDAO.GetAllPaginated(page, pageSize, categoryInt, statusInt, name, sorter)
}
