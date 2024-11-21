package anime

import (
	animesrv "kong-anime-go/internal/services/anime"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	AnimeService *animesrv.AnimeService
}

func NewAnimeHandler(animeService *animesrv.AnimeService) *Handler {
	return &Handler{
		AnimeService: animeService,
	}
}

func (api *Handler) GetTest(c *gin.Context) {
	msg := api.AnimeService.GetTest()
	c.JSON(http.StatusOK, gin.H{"msg": msg})
}
