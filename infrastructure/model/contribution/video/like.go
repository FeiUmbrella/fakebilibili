package video

import "gorm.io/gorm"

type Likes struct {
	gorm.Model
	Uid     uint `json:"uid" gorm:"column:uid"`
	VideoID uint `json:"video_id"  gorm:"column:video_id"`
}

type LikesList []Likes

func (Likes) TableName() string {
	return "lv_video_contribution_like"
}
