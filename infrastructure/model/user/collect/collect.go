package collect

import (
	"fakebilibili/infrastructure/model/contribution/video"
	user2 "fakebilibili/infrastructure/model/user"
	"gorm.io/gorm"
)

// Collect 某个收藏夹与被收藏视频、已经所属用户的关联信息
type Collect struct {
	gorm.Model
	Uid         uint `json:"uid"`
	FavoritesID uint `json:"favorites_id" gorm:"column:favorites_id"`
	VideoID     uint `json:"video_id" gorm:"column:video_id"`

	UserInfo  user2.User               `json:"userInfo" gorm:"foreignKey:Uid"`
	VideoInfo video.VideosContribution `json:"videoInfo" gorm:"foreignKey:VideoID"`
}

type CollectsList []Collect

func (Collect) TableName() string {
	return "lv_users_collect"
}
