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
	VideoID    uint   `json:"video_id"  gorm:"column:video_id"`                   // 进行转码任务的对应视频ID
	Resolution int    `json:"resolution"  gorm:"column:resolution"`               // 分辨率
	Dst        string `json:"dst" gorm:"column:dst;type:varchar(255)"`            // 转码后的保存url
	Status     int    `json:"method"  gorm:"column:status"`                       // 转码任务的状态
	Type       string `json:"path" gorm:"column:type;type:varchar(255)"`          // aliyun
}

func (TranscodingTask) TableName() string {
	return "lv_transcoding_task"
}

// AddTask 保存转码任务信息到数据库
func (t *TranscodingTask) AddTask() bool {
	err := global.MysqlDb.Create(t).Error
	return err == nil
}

// GetInfoByTaskID 根据转码任务ID获取转码详情
func (t *TranscodingTask) GetInfoByTaskID(id string) error {
	return global.MysqlDb.Where("task_id = ?", id).First(t).Error
}

// Save 保存更新后的转码任务详情
func (t *TranscodingTask) Save() bool {
	err := global.MysqlDb.Save(t).Error
	return err == nil
}
