package models

import (
	"kong-anime-go/internal/common"
	"time"

	"gorm.io/gorm"
)

// Follow 追番模型
type Follow struct {
	gorm.Model
	AnimeID    uint                  // 关联的动漫ID
	Anime      Anime                 // 关联的动漫
	Category   common.FollowCategory // 分类 (经典、高质量、新番、厕纸、神作)
	Status     common.FollowStatus   // 状态 (想看、在看、看过)
	FinishedAt *time.Time            `gorm:"default:null"` // 看完时间
}
