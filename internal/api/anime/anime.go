package anime

import (
	"errors"
	"net/http"
	"strings"
	"time"

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
	formats := []string{"200601", "2006-01", "20061", "2006-1"}
	for _, format := range formats {
		t, err := time.Parse(format, season)
		if err == nil {
			return t.Format("2006-01"), nil
		}
	}
	return "", errors.New("invalid season format")
}

func (api *Handler) bindAndValidateAnime(c *gin.Context, anime *models.Anime) ([]string, []string, error) {
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
		return nil, nil, err
	}

	if req.Name == "" || req.Season == "" {
		return nil, nil, errors.New("Name and Season are required")
	}

	formattedSeason, err := formatSeason(req.Season)
	if err != nil {
		return nil, nil, errors.New("Invalid season format")
	}

	anime.Name = req.Name
	anime.Aliases = strings.Join(req.Aliases, ",")
	anime.Production = req.Production
	anime.Season = formattedSeason
	anime.Episodes = req.Episodes
	anime.Image = req.Image

	if len(anime.Image) == 0 { // 如果没有传入图片，则使用 fakeimg.pl 生成图片
		// 计算名字需要的图片尺寸
		l := len(anime.Name)
		w, h := 160, 100
		if l > 10 {
			w = 160 + 10*(l-10)
		}
		anime.Image = "https://fakeimg.pl/" + strconv.Itoa(w) + "x" + strconv.Itoa(h) + "/?text=" + anime.Name + "&font=noto"
	}

	return req.Categories, req.Tags, nil
}

func (api *Handler) CreateAnime(c *gin.Context) {
	anime := &models.Anime{}
	categories, tags, err := api.bindAndValidateAnime(c, anime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anime, err = api.AnimeSrv.CreateAnime(anime, categories, tags)
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

	anime := &models.Anime{}
	anime.ID = uint(id)
	categories, tags, err := api.bindAndValidateAnime(c, anime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anime, err = api.AnimeSrv.UpdateAnime(anime, categories, tags)
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
