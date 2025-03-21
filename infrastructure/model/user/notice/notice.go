package notice

import (
	user2 "fakebilibili/infrastructure/model/user"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// todo:这个表是干啥的？
type Notice struct {
	gorm.Model
	Uid     uint   `json:"uid"`
	Cid     uint   `json:"cid"`
	Type    string `json:"type" gorm:"type:varchar(255)"`
	ToID    uint   `json:"to_id" gorm:"column:to_id"`
	ISRead  uint   `json:"is_read" gorm:"column:is_read"`
	Content string `json:"content" gorm:"type:text"`

	VideoInfo   VideoInfo  `json:"videoInfo" gorm:"foreignKey:to_id"`
	UserInfo    user2.User `json:"userInfo" gorm:"foreignKey:cid"`
	ArticleInfo Article    `json:"articleInfo" gorm:"foreignKey:to_id"`
}

type NoticesList []Notice

func (Notice) TableName() string {
	return "lv_users_notices"
}

// todo:为什么临时加一个VideoInfo能解决依赖循环？
type VideoInfo struct {
	gorm.Model
	Uid   uint           `json:"uid"`
	Title string         `json:"title" gorm:"type:varchar(255)"`
	Video datatypes.JSON `json:"video"`
	Cover datatypes.JSON `json:"cover"`
}

func (VideoInfo) TableName() string {
	return "lv_video_contribution"
}

type Article struct {
	gorm.Model
	Uid              uint           `json:"uid"`
	ClassificationID uint           `json:"classification_id"  gorm:"classification_id"`
	Title            string         `json:"title"`
	Cover            datatypes.JSON `json:"cover"`
}

func (Article) TableName() string {
	return "lv_article_contribution"
}
