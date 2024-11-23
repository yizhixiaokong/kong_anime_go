package anime

import (
	"errors"
	"kong-anime-go/internal/dao"
	"kong-anime-go/internal/dao/models"
	"strings"
)

// Service 处理动漫相关的服务
type Service struct {
	animeDAO    *dao.AnimeDAO
	categoryDAO *dao.CategoryDAO
	tagDAO      *dao.TagDAO
}

// NewService 创建一个新的 AnimeService
func NewService(animeDAO *dao.AnimeDAO, categoryDAO *dao.CategoryDAO, tagDAO *dao.TagDAO) *Service {
	return &Service{
		animeDAO:    animeDAO,
		categoryDAO: categoryDAO,
		tagDAO:      tagDAO,
	}
}

// Create 创建一个新的动漫
func (s *Service) Create(anime *models.Anime, categories []string, tags []string) (*models.Anime, error) {
	if err := s.animeDAO.Create(anime); err != nil {
		return nil, err
	}
	if err := s.addCategoriesToAnime(anime, categories); err != nil {
		return nil, err
	}
	if err := s.addTagsToAnime(anime, tags); err != nil {
		return nil, err
	}
	return s.animeDAO.GetByID(anime.ID)
}

// Delete 删除一个动漫
func (s *Service) Delete(id uint) (uint, error) {
	if err := s.animeDAO.ClearCategories(id); err != nil {
		return 0, err
	}
	if err := s.animeDAO.ClearTags(id); err != nil {
		return 0, err
	}
	// 硬删除
	return id, s.animeDAO.HardDelete(id)
}

// Update 更新一个动漫
func (s *Service) Update(anime *models.Anime, categories []string, tags []string) (*models.Anime, error) {
	existingAnime, err := s.animeDAO.GetByID(anime.ID)
	if err != nil {
		return nil, err
	}
	if existingAnime == nil {
		return nil, errors.New("anime not found")
	}

	existingAnime.Name = anime.Name
	existingAnime.Aliases = anime.Aliases
	existingAnime.Production = anime.Production
	existingAnime.Season = anime.Season
	existingAnime.Episodes = anime.Episodes
	existingAnime.Image = anime.Image

	if err := s.updateCategories(existingAnime, categories); err != nil {
		return nil, err
	}
	if err := s.updateTags(existingAnime, tags); err != nil {
		return nil, err
	}

	if err := s.animeDAO.Update(existingAnime); err != nil {
		return nil, err
	}

	return s.animeDAO.GetByID(existingAnime.ID)
}

func (s *Service) updateCategories(anime *models.Anime, categories []string) error {
	if err := s.animeDAO.ClearCategories(anime.ID); err != nil {
		return err
	}
	return s.addCategoriesToAnime(anime, categories)
}

func (s *Service) updateTags(anime *models.Anime, tags []string) error {
	if err := s.animeDAO.ClearTags(anime.ID); err != nil {
		return err
	}
	return s.addTagsToAnime(anime, tags)
}

func (s *Service) addCategoriesToAnime(anime *models.Anime, categories []string) error {
	for _, categoryName := range categories {
		category, err := s.getOrCreateCategory(categoryName)
		if err != nil {
			return err
		}
		if err := s.animeDAO.AddCategory(anime.ID, category.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) addTagsToAnime(anime *models.Anime, tags []string) error {
	for _, tagName := range tags {
		tag, err := s.getOrCreateTag(tagName)
		if err != nil {
			return err
		}
		if err := s.animeDAO.AddTag(anime.ID, tag.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) getOrCreateCategory(name string) (*models.Category, error) {
	category, err := s.categoryDAO.GetByName(name)
	if err != nil {
		category = &models.Category{Name: name}
		if err := s.categoryDAO.Create(category); err != nil {
			return nil, err
		}
	}
	return category, nil
}

func (s *Service) getOrCreateTag(name string) (*models.Tag, error) {
	tag, err := s.tagDAO.GetByName(name)
	if err != nil {
		tag = &models.Tag{Name: name}
		if err := s.tagDAO.Create(tag); err != nil {
			return nil, err
		}
	}
	return tag, nil
}

// GetByID 根据ID获取动漫
func (s *Service) GetByID(id uint) (*models.Anime, error) {
	return s.animeDAO.GetByID(id)
}

// GetAll 获取所有动漫
func (s *Service) GetAll(page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetAllPaginated(page, pageSize)
}

// GetByName 根据名称获取动漫
func (s *Service) GetByName(name string, page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetByNameAndAlias(name, page, pageSize)
}

// GetBySeason 根据季节获取动漫
func (s *Service) GetBySeason(season string, page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetBySeason(season, page, pageSize)
}

// GetByCategory 根据分类获取动漫
func (s *Service) GetByCategory(categoryName string, page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetByCategory(categoryName, page, pageSize)
}

// GetByTag 根据标签获取动漫
func (s *Service) GetByTag(tagName string, page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetByTag(tagName, page, pageSize)
}

// AddCategoriesToAnime 添加分类到动漫
func (s *Service) AddCategoriesToAnime(animeID uint, categories []string) (*models.Anime, error) {
	anime, err := s.animeDAO.GetByID(animeID)
	if err != nil {
		return nil, err
	}
	if err := s.updateCategories(anime, categories); err != nil {
		return nil, err
	}
	return s.animeDAO.GetByID(anime.ID)
}

// AddTagsToAnime 添加标签到动漫
func (s *Service) AddTagsToAnime(animeID uint, tags []string) (*models.Anime, error) {
	anime, err := s.animeDAO.GetByID(animeID)
	if err != nil {
		return nil, err
	}
	if err := s.updateTags(anime, tags); err != nil {
		return nil, err
	}
	return s.animeDAO.GetByID(anime.ID)
}

// GetAllSeasons 获取所有季节
func (s *Service) GetAllSeasons() (map[string]map[string]int, error) {
	seasons, err := s.animeDAO.GetAllSeasons()
	if err != nil {
		return nil, err
	}

	seasonMap := make(map[string]map[string]int)
	for _, season := range seasons {
		parts := strings.Split(season.Season, "-")
		if len(parts) == 2 {
			year := parts[0]
			month := parts[1]
			if seasonMap[year] == nil {
				seasonMap[year] = make(map[string]int)
			}
			seasonMap[year][month] = season.Count
		}
	}
	return seasonMap, nil
}
