package ping

import (
	pingsrv "kong-anime-go/internal/services/ping"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	PingService *pingsrv.PingService
}

func NewPingHandler(pingservice *pingsrv.PingService) *Handler {
	return &Handler{
		PingService: pingservice,
	}
}

func (api *Handler) GetHello(c *gin.Context) {
	h := api.PingService.GetHello()
	c.JSON(http.StatusOK, gin.H{"msg": h})
}

func (api *Handler) GetPing(c *gin.Context) {
	p, err := api.PingService.GetPing()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, p)
}
