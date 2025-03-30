package checkIn

import (
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
)

// CheckIn 用户签到信息
type CheckIn struct {
	gorm.Model
	Uid             uint `json:"uid"`
	LatestDay       int  `json:"latest_day" gorm:"column:latest_day"`           // 最后一次签到日期
	ConsecutiveDays int  `json:"consecutive_day" gorm:"column:consecutive_day"` // 连续签到天数
	Integral        int  `json:"integral" gorm:"column:integral"`               // 签到获得的积分
}

func (CheckIn) TableName() string {
	return "lv_check_in"
}

// Query 查询签到记录
func (ck *CheckIn) Query() error {
	err := global.MysqlDb.Where("uid=?", ck.Uid).First(ck).Error
	return err
}

// Create 创建一条签到记录
func (ck *CheckIn) Create() bool {
	err := global.MysqlDb.Create(ck).Error
	return err == nil
}

// Update 更新用户记录
func (ck *CheckIn) Update(fields map[string]interface{}) error {
	err := global.MysqlDb.Model(&CheckIn{}).
		Where("uid=?", ck.Uid).
		Updates(fields).Error
	return err
}
