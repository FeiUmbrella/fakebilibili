package checkIn

import "gorm.io/gorm"

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
