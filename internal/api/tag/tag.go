package tag

import (
	"kong-anime-go/internal/dao/models"
	"kong-anime-go/internal/services/tag"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService *tag.TagService
}

func NewTagHandler(tagService *tag.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var newTag models.Tag
	if err := c.ShouldBindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := h.tagService.CreateTag(&newTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	deletedID, err := h.tagService.DeleteTag(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Tag deleted successfully!", "id": deletedID})
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedTag models.Tag
	if err := c.ShouldBindJSON(&updatedTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedTag.ID = uint(id)
	tag, err := h.tagService.UpdateTag(&updatedTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) GetTagByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tag, err := h.tagService.GetTagByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) GetAllTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags, "total": len(tags)})
}

func (h *TagHandler) GetTagsByName(c *gin.Context) {
	name := c.Query("name")
	tags, err := h.tagService.GetTagsByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags, "total": len(tags)})
}

func (h *TagHandler) GetTagStats(c *gin.Context) {
	stats, err := h.tagService.GetTagStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_stats": stats})
}
