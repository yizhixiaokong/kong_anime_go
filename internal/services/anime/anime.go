package anime

import (
	"errors"
	"kong-anime-go/internal/dao"
	"kong-anime-go/internal/dao/models"
)

type AnimeService struct {
	animeDAO    *dao.AnimeDAO
	categoryDAO *dao.CategoryDAO
	tagDAO      *dao.TagDAO
}

func NewAnimeService(animeDAO *dao.AnimeDAO, categoryDAO *dao.CategoryDAO, tagDAO *dao.TagDAO) *AnimeService {
	return &AnimeService{
		animeDAO:    animeDAO,
		categoryDAO: categoryDAO,
		tagDAO:      tagDAO,
	}
}

func (s *AnimeService) CreateAnime(anime *models.Anime, categories []string, tags []string) (*models.Anime, error) {
	for _, categoryName := range categories {
		category, err := s.categoryDAO.GetByName(categoryName)
		if err != nil {
			category = &models.Category{Name: categoryName}
			if err := s.categoryDAO.Create(category); err != nil {
				return nil, err
			}
		}
		anime.Categories = append(anime.Categories, *category)
	}

	for _, tagName := range tags {
		tag, err := s.tagDAO.GetByName(tagName)
		if err != nil {
			tag = &models.Tag{Name: tagName}
			if err := s.tagDAO.Create(tag); err != nil {
				return nil, err
			}
		}
		anime.Tags = append(anime.Tags, *tag)
	}

	return anime, s.animeDAO.Create(anime)
}

func (s *AnimeService) DeleteAnime(id uint) (uint, error) {
	return id, s.animeDAO.Delete(id)
}

func (s *AnimeService) UpdateAnime(anime *models.Anime, categories []string, tags []string) (*models.Anime, error) {
	existingAnime, err := s.animeDAO.GetByID(anime.ID)
	if err != nil {
		return nil, err
	}
	if existingAnime == nil {
		return nil, errors.New("anime not found")
	}

	anime.Categories = nil
	anime.Tags = nil

	for _, categoryName := range categories {
		category, err := s.categoryDAO.GetByName(categoryName)
		if err != nil {
			category = &models.Category{Name: categoryName}
			if err := s.categoryDAO.Create(category); err != nil {
				return nil, err
			}
		}
		anime.Categories = append(anime.Categories, *category)
	}

	for _, tagName := range tags {
		tag, err := s.tagDAO.GetByName(tagName)
		if err != nil {
			tag = &models.Tag{Name: tagName}
			if err := s.tagDAO.Create(tag); err != nil {
				return nil, err
			}
		}
		anime.Tags = append(anime.Tags, *tag)
	}

	return anime, s.animeDAO.Update(anime)
}

func (s *AnimeService) GetAnimeByID(id uint) (*models.Anime, error) {
	return s.animeDAO.GetByID(id)
}

func (s *AnimeService) GetAllAnimes(page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetAllPaginated(page, pageSize)
}

func (s *AnimeService) GetAnimesByName(name string, page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetByNameAndAlias(name, page, pageSize)
}

func (s *AnimeService) GetAnimesBySeason(season string, page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetBySeason(season, page, pageSize)
}

func (s *AnimeService) GetAnimesByCategory(categoryName string, page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetByCategory(categoryName, page, pageSize)
}

func (s *AnimeService) GetAnimesByTag(tagName string, page, pageSize int) ([]models.Anime, int64, error) {
	return s.animeDAO.GetByTag(tagName, page, pageSize)
}

func (s *AnimeService) AddCategoriesToAnime(animeID uint, categories []string) (*models.Anime, error) {
	anime, err := s.animeDAO.GetByID(animeID)
	if err != nil {
		return nil, err
	}

	for _, categoryName := range categories {
		category, err := s.categoryDAO.GetByName(categoryName)
		if err != nil {
			category = &models.Category{Name: categoryName}
			if err := s.categoryDAO.Create(category); err != nil {
				return nil, err
			}
		}
		// Check if category already exists in anime
		exists := false
		for _, c := range anime.Categories {
			if c.ID == category.ID {
				exists = true
				break
			}
		}
		if !exists {
			anime.Categories = append(anime.Categories, *category)
		}
	}

	return anime, s.animeDAO.Update(anime)
}

func (s *AnimeService) AddTagsToAnime(animeID uint, tags []string) (*models.Anime, error) {
	anime, err := s.animeDAO.GetByID(animeID)
	if err != nil {
		return nil, err
	}

	for _, tagName := range tags {
		tag, err := s.tagDAO.GetByName(tagName)
		if err != nil {
			tag = &models.Tag{Name: tagName}
			if err := s.tagDAO.Create(tag); err != nil {
				return nil, err
			}
		}
		// Check if tag already exists in anime
		exists := false
		for _, t := range anime.Tags {
			if t.ID == tag.ID {
				exists = true
				break
			}
		}
		if !exists {
			anime.Tags = append(anime.Tags, *tag)
		}
	}

	return anime, s.animeDAO.Update(anime)
}
