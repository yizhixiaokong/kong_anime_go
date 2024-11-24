package follow

import (
	"net/http"
	"strconv"
	"time"

	"kong-anime-go/internal/common"
	"kong-anime-go/internal/dao/models"
	"kong-anime-go/internal/services/follow"

	"github.com/gin-gonic/gin"
)

// Handler 处理追番相关的HTTP请求
type Handler struct {
	service *follow.Service
}

// NewHandler 创建一个新的 FollowHandler
func NewHandler(service *follow.Service) *Handler {
	return &Handler{service: service}
}

// Create 创建追番
func (h *Handler) Create(c *gin.Context) {
	var follow models.Follow
	if err := c.ShouldBindJSON(&follow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !follow.Category.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
		return
	}
	if !follow.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	existingFollow, err := h.service.GetByAnimeID(follow.AnimeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingFollow != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Follow already exists for this anime"})
		return
	}
	follow.FinishedAt = nil // 设置为空值

	createdFollow, err := h.service.Create(&follow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdFollow)
}

// Delete 删除追番
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if _, err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully", "id": id})
}

// Update 更新追番
func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedFollow models.Follow
	if err := c.ShouldBindJSON(&updatedFollow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !updatedFollow.Category.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
		return
	}
	if !updatedFollow.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	updatedFollow.ID = uint(id)
	follow, err := h.service.Update(&updatedFollow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, follow)
}

// GetByID 根据ID获取追番
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	follow, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, follow)
}

// GetAll 获取所有追番
func (h *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	categoryStr := c.DefaultQuery("category", "")
	statusStr := c.DefaultQuery("status", "")
	var category, status *int
	if categoryStr != "" {
		categoryVal, _ := strconv.Atoi(categoryStr)
		category = &categoryVal
	}
	if statusStr != "" {
		statusVal, _ := strconv.Atoi(statusStr)
		status = &statusVal
	}
	follows, total, err := h.service.GetAll(page, pageSize, category, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": follows, "total": total})
}

// UpdateStatus 更新追番状态
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var request struct {
		Status common.FollowStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !request.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	follow, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	follow.Status = request.Status
	if request.Status == common.FollowStatusWatched {
		now := time.Now()
		follow.FinishedAt = &now
	} else {
		follow.FinishedAt = nil // 设置为空值
	}
	updatedFollow, err := h.service.Update(follow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedFollow)
}

// GetAllCategories 获取所有追番分类
func (h *Handler) GetAllCategories(c *gin.Context) {
	type CategoryResponse struct {
		Value       int    `json:"value"`
		String      string `json:"string"`
		Description string `json:"description"`
	}

	var categories []CategoryResponse
	for _, category := range common.AllFollowCategories() {
		categories = append(categories, CategoryResponse{
			Value:       int(category),
			String:      category.String(),
			Description: category.Description(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
