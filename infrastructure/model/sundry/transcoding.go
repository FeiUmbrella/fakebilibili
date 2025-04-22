package sundry

import (
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
)

var Aliyun = "aliyun"

// TranscodingTask 记录使用进行的转码任务
type TranscodingTask struct {
	gorm.Model
	TaskID     string `json:"interface"  gorm:"column:task_id;type:varchar(255)"` // 任务id
	VideoID    uint   `json:"video_id"  gorm:"column:video_id"`
	Resolution int    `json:"resolution"  gorm:"column:resolution"`    // 分辨率
	Dst        string `json:"dst" gorm:"column:dst;type:varchar(255)"` // 转码后的保存url
	Status     int    `json:"method"  gorm:"column:status"`
	Type       string `json:"path" gorm:"column:type;type:varchar(255)"` // aliyun
}

func (TranscodingTask) TableName() string {
	return "lv_transcoding_task"
}

// AddTask 保存转码任务信息到数据库
func (t *TranscodingTask) AddTask() bool {
	err := global.MysqlDb.Create(t).Error
	return err == nil
}
