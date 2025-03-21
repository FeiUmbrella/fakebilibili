package record

import (
	"fakebilibili/infrastructure/model/contribution/article"
	"fakebilibili/infrastructure/model/contribution/video"
	user2 "fakebilibili/infrastructure/model/user"
	"gorm.io/gorm"
)

// Record 浏览记录
type Record struct {
	gorm.Model
	Uid  uint   `json:"uid"`
	Type string `json:"type" gorm:"type:varchar(255)"` // todo:这个字段什么用？
	ToId uint   `json:"to_id" gorm:"column:to_id"`     // todo:这个字段什么用？

	VideoInfo   video.VideosContribution     `json:"videoInfo" gorm:"foreignKey:to_id"`
	UserInfo    user2.User                   `json:"userInfo" gorm:"foreignKey:uid"`
	ArticleInfo article.ArticlesContribution `json:"articleInfo" gorm:"foreignKey:to_id"`
}
type RecordList []Record

func (Record) TableName() string {
	return "lv_users_record"
}
