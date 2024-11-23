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

// Handler 处理 Anime 相关的服务
type Handler struct {
	AnimeSrv *animesrv.Service
}

// NewHandler 创建一个新的 AnimeHandler
func NewHandler(animeSrv *animesrv.Service) *Handler {
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

// Create 创建一个新的动漫
func (api *Handler) Create(c *gin.Context) {
	anime := &models.Anime{}
	categories, tags, err := api.bindAndValidateAnime(c, anime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anime, err = api.AnimeSrv.Create(anime, categories, tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Anime created successfully!", "anime": anime})
}

// Delete 删除一个动漫
func (api *Handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if _, err := api.AnimeSrv.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "id": id})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Anime deleted successfully!", "id": id})
}

// Update 更新一个动漫
func (api *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	anime := &models.Anime{}
	anime.ID = uint(id)
	categories, tags, err := api.bindAndValidateAnime(c, anime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anime, err = api.AnimeSrv.Update(anime, categories, tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Anime updated successfully!", "anime": anime})
}

// GetByID 根据ID获取动漫
func (api *Handler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	anime, err := api.AnimeSrv.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"anime": anime})
}

// GetAll 获取所有动漫
func (api *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

// GetByName 根据名称获取动漫
func (api *Handler) GetByName(c *gin.Context) {
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetByName(name, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

// GetBySeason 根据季节获取动漫
func (api *Handler) GetBySeason(c *gin.Context) {
	season := c.Query("season")
	formattedSeason, err := formatSeason(season)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season format"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetBySeason(formattedSeason, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

// GetByCategory 根据分类获取动漫
func (api *Handler) GetByCategory(c *gin.Context) {
	categoryName := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetByCategory(categoryName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

// GetByTag 根据标签获取动漫
func (api *Handler) GetByTag(c *gin.Context) {
	tagName := c.Query("tag")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	animes, total, err := api.AnimeSrv.GetByTag(tagName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"animes": animes, "total": total, "page": page, "pageSize": pageSize})
}

// AddCategoriesToAnime 添加分类到动漫
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

// AddTagsToAnime 添加标签到动漫
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

// GetAllSeasons 获取所有季节
func (api *Handler) GetAllSeasons(c *gin.Context) {
	seasons, err := api.AnimeSrv.GetAllSeasons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seasons": seasons})
}
