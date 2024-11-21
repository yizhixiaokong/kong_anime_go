package ping

import (
	"kong-anime-go/internal/dao/models"
	"time"
)

type PingService struct{}

func NewPingService() *PingService {
	return &PingService{}
}

func (s *PingService) GetHello() string {
	return "Hello, World!"
}

func (s *PingService) GetPing() (*models.Ping, error) {
	return &models.Ping{
		Msg:  "pong",
		Time: time.Now().Format(time.DateTime),
	}, nil
}
