package ping

import (
	pingsrv "kong-anime-go/internal/services/ping"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 处理 Ping 相关的服务
type Handler struct {
	PingService *pingsrv.Service
}

// NewHandler 创建一个新的 PingHandler
func NewHandler(pingservice *pingsrv.Service) *Handler {
	return &Handler{
		PingService: pingservice,
	}
}

// GetHello 返回一个问候消息
func (api *Handler) GetHello(c *gin.Context) {
	h := api.PingService.GetHello()
	c.JSON(http.StatusOK, gin.H{"msg": h})
}

// GetPing 返回一个 Ping 消息
func (api *Handler) GetPing(c *gin.Context) {
	p, err := api.PingService.GetPing()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, p)
}
