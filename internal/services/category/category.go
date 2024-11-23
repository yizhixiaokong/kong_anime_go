package category

import (
	"errors"
	"kong-anime-go/internal/dao"
	"kong-anime-go/internal/dao/models"
)

type CategoryService struct {
	categoryDAO *dao.CategoryDAO
}

func NewCategoryService(categoryDAO *dao.CategoryDAO) *CategoryService {
	return &CategoryService{
		categoryDAO: categoryDAO,
	}
}

func (s *CategoryService) CreateCategory(category *models.Category) (*models.Category, error) {
	existingCategory, err := s.categoryDAO.GetByName(category.Name)
	if err == nil && existingCategory != nil {
		return existingCategory, nil
	}
	if err := s.categoryDAO.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(id uint) (uint, error) {
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

func (s *CategoryService) UpdateCategory(category *models.Category) (*models.Category, error) {
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

func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	return s.categoryDAO.GetByID(id)
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryDAO.GetAll()
}

func (s *CategoryService) GetCategoriesByName(name string) ([]models.Category, error) {
	return s.categoryDAO.GetByNameLike(name)
}

func (s *CategoryService) GetCategoryStats() (map[string]int, error) {
	return s.categoryDAO.GetCategoryStats()
}
