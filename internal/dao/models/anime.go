package models

import (
	"gorm.io/gorm"
)

// Anime 动漫模型
type Anime struct {
	gorm.Model
	Name       string     `gorm:"not null"`                    // 名称
	Aliases    string     `gorm:"type:text"`                   // 别名(原名、昵称等)，用逗号隔开
	Categories []Category `gorm:"many2many:anime_categories;"` // 分类 (热血、冒险、搞笑、奇幻等)
	Tags       []Tag      `gorm:"many2many:anime_tags;"`       // 标签 (原创、漫改、游戏改、小说改、其他)
	Production string     // 制作公司
	Season     string     // 季度(包含年份)
	Episodes   int        // 集数
	Image      string     // 图片（可能存储为Base64）
}
