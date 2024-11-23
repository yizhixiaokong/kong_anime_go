package tag

import (
	"kong-anime-go/internal/dao/models"
	"kong-anime-go/internal/services/tag"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 处理 Tag 相关的服务
type Handler struct {
	tagService *tag.Service
}

// NewHandler 创建一个新的 TagHandler
func NewHandler(tagService *tag.Service) *Handler {
	return &Handler{
		tagService: tagService,
	}
}

// Create 创建一个新的标签
func (h *Handler) Create(c *gin.Context) {
	var newTag models.Tag
	if err := c.ShouldBindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := h.tagService.Create(&newTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// Delete 删除一个标签
func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	deletedID, err := h.tagService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Tag deleted successfully!", "id": deletedID})
}

// Update 更新一个标签
func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedTag models.Tag
	if err := c.ShouldBindJSON(&updatedTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedTag.ID = uint(id)
	tag, err := h.tagService.Update(&updatedTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// GetByID 根据ID获取标签
func (h *Handler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tag, err := h.tagService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// GetAll 获取所有标签
func (h *Handler) GetAll(c *gin.Context) {
	tags, err := h.tagService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags, "total": len(tags)})
}

// GetByName 根据名称获取标签
func (h *Handler) GetByName(c *gin.Context) {
	name := c.Query("name")
	tags, err := h.tagService.GetByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags, "total": len(tags)})
}

// GetStats 获取标签统计信息
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.tagService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_stats": stats})
}
