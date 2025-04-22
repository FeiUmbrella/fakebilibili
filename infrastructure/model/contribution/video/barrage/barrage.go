package barrage

import (
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Barrage 弹幕
type Barrage struct {
	gorm.Model
	Uid     uint    `json:"uid"`
	VideoID uint    `json:"video_id" gorm:"video_id"`
	Time    float64 `json:"time"`
	Author  string  `json:"author" gorm:"type:varchar(255)"`
	Type    uint    `json:"type"`
	Text    string  `json:"text" gorm:"type:text"`
	Color   uint    `json:"color"`

	UserInfo  user.User `json:"user_info" gorm:"foreignKey:Uid"`
	VideoInfo VideoInfo `json:"video_info" gorm:"foreignKey:VideoID"`
}
type BarragesList []Barrage

func (Barrage) TableName() string {
	return "lv_video_contribution_barrage"
}

// VideoInfo 临时加一个video模型解决依赖循环
type VideoInfo struct {
	gorm.Model
	Uid   uint           `json:"uid" gorm:"uid"`
	Title string         `json:"title" gorm:"title"`
	Video datatypes.JSON `json:"video" gorm:"video"`
	Cover datatypes.JSON `json:"cover" gorm:"cover"`
}

func (VideoInfo) TableName() string {
	return "lv_video_contribution"
}

// GetVideoBarrageByID 查询视频弹幕
func (bl *BarragesList) GetVideoBarrageByID(vid uint) bool {
	err := global.MysqlDb.Where("video_id = ?", vid).Find(&bl).Error
	return err == nil
}

// Create 创建弹幕
func (b *Barrage) Create() bool {
	err := global.MysqlDb.Create(&b).Error
	return err == nil
}
