package anime

import (
	"errors"
	"kong-anime-go/internal/dao"
	"kong-anime-go/internal/dao/models"
	"strings"
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

func (s *AnimeService) DeleteAnime(id uint) (uint, error) {
	if err := s.animeDAO.ClearCategories(id); err != nil {
		return 0, err
	}
	if err := s.animeDAO.ClearTags(id); err != nil {
		return 0, err
	}
	// 硬删除
	return id, s.animeDAO.HardDelete(id)
}

func (s *AnimeService) UpdateAnime(anime *models.Anime, categories []string, tags []string) (*models.Anime, error) {
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

func (s *AnimeService) updateCategories(anime *models.Anime, categories []string) error {
	if err := s.animeDAO.ClearCategories(anime.ID); err != nil {
		return err
	}
	return s.addCategoriesToAnime(anime, categories)
}

func (s *AnimeService) updateTags(anime *models.Anime, tags []string) error {
	if err := s.animeDAO.ClearTags(anime.ID); err != nil {
		return err
	}
	return s.addTagsToAnime(anime, tags)
}

func (s *AnimeService) addCategoriesToAnime(anime *models.Anime, categories []string) error {
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

func (s *AnimeService) addTagsToAnime(anime *models.Anime, tags []string) error {
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

func (s *AnimeService) getOrCreateCategory(name string) (*models.Category, error) {
	category, err := s.categoryDAO.GetByName(name)
	if err != nil {
		category = &models.Category{Name: name}
		if err := s.categoryDAO.Create(category); err != nil {
			return nil, err
		}
	}
	return category, nil
}

func (s *AnimeService) getOrCreateTag(name string) (*models.Tag, error) {
	tag, err := s.tagDAO.GetByName(name)
	if err != nil {
		tag = &models.Tag{Name: name}
		if err := s.tagDAO.Create(tag); err != nil {
			return nil, err
		}
	}
	return tag, nil
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
	if err := s.updateCategories(anime, categories); err != nil {
		return nil, err
	}
	return s.animeDAO.GetByID(anime.ID)
}

func (s *AnimeService) AddTagsToAnime(animeID uint, tags []string) (*models.Anime, error) {
	anime, err := s.animeDAO.GetByID(animeID)
	if err != nil {
		return nil, err
	}
	if err := s.updateTags(anime, tags); err != nil {
		return nil, err
	}
	return s.animeDAO.GetByID(anime.ID)
}

func (s *AnimeService) GetAllSeasons() (map[string]map[string]int, error) {
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
