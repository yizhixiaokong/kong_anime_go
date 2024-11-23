package category

import (
	"kong-anime-go/internal/dao/models"
	"kong-anime-go/internal/services/category"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 处理 Category 相关的服务
type Handler struct {
	categoryService *category.Service
}

// NewHandler 创建一个新的 CategoryHandler
func NewHandler(categoryService *category.Service) *Handler {
	return &Handler{
		categoryService: categoryService,
	}
}

// Create 创建一个新的分类
func (h *Handler) Create(c *gin.Context) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := h.categoryService.Create(&newCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// Delete 删除一个分类
func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	deletedID, err := h.categoryService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Category deleted successfully!", "id": deletedID})
}

// Update 更新一个分类
func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedCategory.ID = uint(id)
	category, err := h.categoryService.Update(&updatedCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// GetByID 根据ID获取分类
func (h *Handler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// GetAll 获取所有分类
func (h *Handler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories, "total": len(categories)})
}

// GetByName 根据名称获取分类
func (h *Handler) GetByName(c *gin.Context) {
	name := c.Query("name")
	categories, err := h.categoryService.GetByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories, "total": len(categories)})
}

// GetStats 获取分类统计信息
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.categoryService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category_stats": stats})
}
