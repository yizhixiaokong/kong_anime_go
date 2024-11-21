package anime

import (
	"errors"
	"net/http"
	"strings"
	"time" // 添加time包

	"kong-anime-go/internal/dao/models"
	animesrv "kong-anime-go/internal/services/anime"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	AnimeSrv *animesrv.AnimeService
}

func NewAnimeHandler(animeSrv *animesrv.AnimeService) *Handler {
	return &Handler{
		AnimeSrv: animeSrv,
	}
}

func formatSeason(season string) (string, error) {
	// 尝试解析不符合格式的Season
	formats := []string{"200601", "2006-01", "20061", "2006-1"}
	for _, format := range formats {
		t, err := time.Parse(format, season)
		if err == nil {
			// log.Println(t)
			return t.Format("2006-01"), nil
		}
	}
	return "", errors.New("invalid season format")
}

func (api *Handler) CreateAnime(c *gin.Context) {
	var req struct {
		Name       string   `json:"name"`
		Aliases    []string `json:"aliases"`
		Categories []string `json:"categories"`
		Tags       []string `json:"tags"`
		Production string   `json:"production"`
		Season     string   `json:"season"`
		Episodes   int      `json:"episodes"`
		Image      string   `json:"image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" || req.Season == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Season are required"})
		return
	}

	// 格式化Season字段
	formattedSeason, err := formatSeason(req.Season)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season format"})
		return
	}

	anime := &models.Anime{
		Name:       req.Name,
		Aliases:    strings.Join(req.Aliases, ","),
		Production: req.Production,
		Season:     formattedSeason,
		Episodes:   req.Episodes,
		Image:      req.Image,
	}

	anime, err = api.AnimeSrv.CreateAnime(anime, req.Categories, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Anime created successfully!", "anime": anime})
}

func (api *Handler) DeleteAnime(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if _, err := api.AnimeSrv.DeleteAnime(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "id": id})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Anime deleted successfully!", "id": id})
}

func (api *Handler) UpdateAnime(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Name       string   `json:"name"`
		Aliases    []string `json:"aliases"`
		Categories []string `json:"categories"`
		Tags       []string `json:"tags"`
		Production string   `json:"production"`
		Season     string   `json:"season"`
		Episodes   int      `json:"episodes"`
		Image      string   `json:"image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" || req.Season == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Season are required"})
		return
	}

	// 格式化Season字段
	formattedSeason, err := formatSeason(req.Season)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season format"})
		return
	}

	anime, err := api.AnimeSrv.GetAnimeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

	anime.Name = req.Name
	anime.Aliases = strings.Join(req.Aliases, ",")
	anime.Production = req.Production
	anime.Season = formattedSeason
	anime.Episodes = req.Episodes
	anime.Image = req.Image

	anime, err = api.AnimeSrv.UpdateAnime(anime, req.Categories, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Anime updated successfully!", "anime": anime})
}

func (api *Handler) GetAnimeByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	anime, err := api.AnimeSrv.GetAnimeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"anime": anime})
}

func (api *Handler) GetAllAnimes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetAllAnimes(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

func (api *Handler) GetAnimesByName(c *gin.Context) {
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetAnimesByName(name, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

func (api *Handler) GetAnimesBySeason(c *gin.Context) {
	season := c.Query("season")
	// 格式化Season字段
	formattedSeason, err := formatSeason(season)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season format"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetAnimesBySeason(formattedSeason, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

func (api *Handler) GetAnimesByCategory(c *gin.Context) {
	categoryName := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetAnimesByCategory(categoryName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

func (api *Handler) GetAnimesByTag(c *gin.Context) {
	tagName := c.Query("tag")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetAnimesByTag(tagName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

func (api *Handler) AddCategoriesToAnime(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Categories []string `json:"categories"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anime, err := api.AnimeSrv.AddCategoriesToAnime(uint(id), req.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Categories added successfully!", "anime": anime})
}

func (api *Handler) AddTagsToAnime(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Tags []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anime, err := api.AnimeSrv.AddTagsToAnime(uint(id), req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Tags added successfully!", "anime": anime})
}
