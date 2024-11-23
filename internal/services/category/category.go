package category

import (
	"errors"
	"kong-anime-go/internal/dao"
	"kong-anime-go/internal/dao/models"
)

// Service 处理分类相关的服务
type Service struct {
	categoryDAO *dao.CategoryDAO
}

// NewService 创建一个新的 CategoryService
func NewService(categoryDAO *dao.CategoryDAO) *Service {
	return &Service{
		categoryDAO: categoryDAO,
	}
}

// Create 创建一个新的分类
func (s *Service) Create(category *models.Category) (*models.Category, error) {
	existingCategory, err := s.categoryDAO.GetByName(category.Name)
	if err == nil && existingCategory != nil {
		return existingCategory, nil
	}
	if err := s.categoryDAO.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

// Delete 删除一个分类
func (s *Service) Delete(id uint) (uint, error) {
	// 检查是否有相关联的项
	relatedItems, err := s.categoryDAO.CheckRelatedItems(id)
	if err != nil {
		return 0, err
	}
	if relatedItems {
		return 0, errors.New("cannot delete category with related items")
	}
	// 硬删除
	if err := s.categoryDAO.HardDelete(id); err != nil {
		return 0, err
	}
	return id, nil
}

// Update 更新一个分类
func (s *Service) Update(category *models.Category) (*models.Category, error) {
	existingCategory, err := s.categoryDAO.GetByID(category.ID)
	if err != nil {
		return nil, err
	}
	if existingCategory == nil {
		return nil, errors.New("category not found")
	}

	if existingCategory.Name == category.Name {
		return existingCategory, nil
	}

	existingCategory.Name = category.Name
	if err := s.categoryDAO.Update(existingCategory); err != nil {
		return nil, err
	}
	return existingCategory, nil
}

// GetByID 根据ID获取分类
func (s *Service) GetByID(id uint) (*models.Category, error) {
	return s.categoryDAO.GetByID(id)
}

// GetAll 获取所有分类
func (s *Service) GetAll() ([]models.Category, error) {
	return s.categoryDAO.GetAll()
}

// GetByName 根据名称获取分类
func (s *Service) GetByName(name string) ([]models.Category, error) {
	return s.categoryDAO.GetByNameLike(name)
}

// GetStats 获取分类统计信息
func (s *Service) GetStats() (map[string]int, error) {
	return s.categoryDAO.GetCategoryStats()
}
