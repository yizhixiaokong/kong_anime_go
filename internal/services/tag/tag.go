package tag

import (
	"errors"
	"kong-anime-go/internal/dao"
	"kong-anime-go/internal/dao/models"
)

// Service 处理标签相关的服务
type Service struct {
	tagDAO *dao.TagDAO
}

// NewService 创建一个新的 TagService
func NewService(tagDAO *dao.TagDAO) *Service {
	return &Service{
		tagDAO: tagDAO,
	}
}

// Create 创建一个新的标签
func (s *Service) Create(tag *models.Tag) (*models.Tag, error) {
	existingTag, err := s.tagDAO.GetByName(tag.Name)
	if err == nil && existingTag != nil {
		return existingTag, nil
	}
	if err := s.tagDAO.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

// Delete 删除一个标签
func (s *Service) Delete(id uint) (uint, error) {
	// 检查是否有相关联的项
	relatedItems, err := s.tagDAO.CheckRelatedItems(id)
	if err != nil {
		return 0, err
	}
	if relatedItems {
		return 0, errors.New("cannot delete tag with related items")
	}
	// 硬删除
	if err := s.tagDAO.HardDelete(id); err != nil {
		return 0, err
	}
	return id, nil
}

// Update 更新一个标签
func (s *Service) Update(tag *models.Tag) (*models.Tag, error) {
	existingTag, err := s.tagDAO.GetByID(tag.ID)
	if err != nil {
		return nil, err
	}
	if existingTag == nil {
		return nil, errors.New("tag not found")
	}

	if existingTag.Name == tag.Name {
		return existingTag, nil
	}

	existingTag.Name = tag.Name
	if err := s.tagDAO.Update(existingTag); err != nil {
		return nil, err
	}
	return existingTag, nil
}

// GetByID 根据ID获取标签
func (s *Service) GetByID(id uint) (*models.Tag, error) {
	return s.tagDAO.GetByID(id)
}

// GetAll 获取所有标签
func (s *Service) GetAll() ([]models.Tag, error) {
	return s.tagDAO.GetAll()
}

// GetByName 根据名称获取标签
func (s *Service) GetByName(name string) ([]models.Tag, error) {
	return s.tagDAO.GetByNameLike(name)
}

// GetStats 获取标签统计信息
func (s *Service) GetStats() (map[string]int, error) {
	return s.tagDAO.GetTagStats()
}
