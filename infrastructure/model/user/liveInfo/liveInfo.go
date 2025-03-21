package liveInfo

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// LiveInfo 直播间信息
type LiveInfo struct {
	gorm.Model
	Uid   uint           `json:"uid"`                            // 用户id
	Title string         `json:"title" gorm:"type:varchar(255)"` // 直播标题
	Img   datatypes.JSON `json:"img" gorm:"comment:img"`         // 直播封面
}

func (LiveInfo) TableName() string {
	return "lv_live_info"
}
