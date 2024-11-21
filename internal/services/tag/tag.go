package tag

import (
	"errors"
	"kong-anime-go/internal/dao"
	"kong-anime-go/internal/dao/models"
)

type TagService struct {
	tagDAO *dao.TagDAO
}

func NewTagService(tagDAO *dao.TagDAO) *TagService {
	return &TagService{
		tagDAO: tagDAO,
	}
}

func (s *TagService) CreateTag(tag *models.Tag) (*models.Tag, error) {
	existingTag, err := s.tagDAO.GetByName(tag.Name)
	if err == nil && existingTag != nil {
		return existingTag, nil
	}
	if err := s.tagDAO.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) DeleteTag(id uint) (uint, error) {
	// 检查是否有相关联的项
	relatedItems, err := s.tagDAO.CheckRelatedItems(id)
	if err != nil {
		return 0, err
	}
	if relatedItems {
		return 0, errors.New("cannot delete tag with related items")
	}
	if err := s.tagDAO.Delete(id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *TagService) UpdateTag(tag *models.Tag) (*models.Tag, error) {
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

func (s *TagService) GetTagByID(id uint) (*models.Tag, error) {
	return s.tagDAO.GetByID(id)
}

func (s *TagService) GetAllTags() ([]models.Tag, error) {
	return s.tagDAO.GetAll()
}
