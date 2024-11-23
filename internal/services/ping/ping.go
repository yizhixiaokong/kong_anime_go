package ping

import (
	"kong-anime-go/internal/dao/models"
	"time"
)

// Service 处理 Ping 相关的服务
type Service struct{}

// NewService 创建一个新的 PingService
func NewService() *Service {
	return &Service{}
}

// GetHello 返回一个问候消息
func (s *Service) GetHello() string {
	return "Hello, World!"
}

// GetPing 返回一个 Ping 消息
func (s *Service) GetPing() (*models.Ping, error) {
	return &models.Ping{
		Msg:  "pong",
		Time: time.Now().Format(time.DateTime),
	}, nil
}
