package favorites

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Favorites 用户收藏夹信息
type Favorites struct {
	gorm.Model
	Uid     uint           `json:"uid"`                                  // 收藏夹所属用户
	Title   string         `json:"title" gorm:"type:varchar(255)"`       // 收藏夹标题
	Content string         `json:"content" gorm:"type:text"`             // 收藏夹简介
	Cover   datatypes.JSON `json:"cover" gorm:"type:json;comment:cover"` // 收藏夹封面图片链接
	Max     int            `json:"max"`                                  // 单个收藏夹最大收藏视频数
}

type FavoriteList []Favorites

func (Favorites) TableName() string {
	return "lv_users_favorites"
}
