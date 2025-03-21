package video

import (
	"fakebilibili/infrastructure/model/contribution/video/barrage"
	"fakebilibili/infrastructure/model/contribution/video/comments"
	"fakebilibili/infrastructure/model/user"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type VideosContribution struct {
	gorm.Model
	Uid           uint           `json:"uid" gorm:"column:uid"`
	Title         string         `json:"title" gorm:"column:title;type:varchar(255)"`
	Video         datatypes.JSON `json:"video" gorm:"column:video"` //默认1080p
	Video720p     datatypes.JSON `json:"video_720p" gorm:"column:video_720p"`
	Video480p     datatypes.JSON `json:"video_480p" gorm:"column:video_480p"`
	Video360p     datatypes.JSON `json:"video_360p" gorm:"column:video_360p"`
	MediaID       string         `json:"media_id" gorm:"column:media_id; type:varchar(255)"`
	Cover         datatypes.JSON `json:"cover" gorm:"column:cover"`
	VideoDuration int64          `json:"video_duration" gorm:"column:video_duration"`
	Reprinted     int8           `json:"reprinted" gorm:"column:reprinted"`
	Label         string         `json:"label" gorm:"column:label; type:varchar(255)"`
	Introduce     string         `json:"introduce" gorm:"column:introduce; type:varchar(255)"`
	Heat          int            `json:"heat" gorm:"column:heat"`
	//todo:加了一个visible字段，可能会引起很多连锁反应
	IsVisible int `json:"is_visible" gorm:"column:is_visible"`

	UserInfo user.User            `json:"user_info" gorm:"foreignKey:Uid"`
	Likes    LikesList            `json:"likes" gorm:"foreignKey:VideoID" `
	Comments comments.CommentList `json:"comments" gorm:"foreignKey:VideoID"`
	Barrage  barrage.BarragesList `json:"barrage" gorm:"foreignKey:VideoID"`
}

type VideosContributionList []VideosContribution

func (VideosContribution) TableName() string {
	return "lv_video_contribution"
}
