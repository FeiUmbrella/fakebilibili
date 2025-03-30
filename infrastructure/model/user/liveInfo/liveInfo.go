package liveInfo

import (
	"fakebilibili/infrastructure/pkg/global"
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

// IsExistByField 查找是否存在对应字段的记录
func (lI *LiveInfo) IsExistByField(field string, value any) bool {
	err := global.MysqlDb.Model(&LiveInfo{}).Where(field+" = ?", value).First(&lI).Error
	return err == nil
}

// Create 创建直播间
func (lI *LiveInfo) Create() bool {
	err := global.MysqlDb.
		Model(&LiveInfo{}).
		Create(&lI).Error
	return err == nil
}

// UpdateInfo 修改直播间信息
func (lI *LiveInfo) UpdateInfo() bool {
	t := new(LiveInfo)
	if t.IsExistByField("uid", lI.Uid) {
		// 存在该直播间
		err := global.MysqlDb.Model(&LiveInfo{}).
			Where("uid = ?", lI.Uid).
			Updates(&lI).Error
		return err == nil
	}

	// 不存在，创建该直播间
	return lI.Create()
}
