package cron_events

import "gorm.io/gorm"

// RuntimeLogEntry 日志的结构体
type RuntimeLogEntry struct {
	gorm.Model
	Time     string `json:"time;type:varchar(255)"`
	Level    string `json:"level;type:varchar(255)"`
	Msg      string `json:"msg;type:varchar(255)"`
	File     string `json:"file;type:varchar(255)"` //避免解析info日志出现null值
	Function string `json:"function;type:varchar(255)"`
}

func (RuntimeLogEntry) TableName() string {
	return "lv_runtime_log"
}
