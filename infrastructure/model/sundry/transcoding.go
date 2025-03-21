package sundry

import "gorm.io/gorm"

// todo:这个表干啥使得？
type TranscodingTask struct {
	gorm.Model
	TaskID     string `json:"interface"  gorm:"column:task_id;type:varchar(255)"`
	VideoID    uint   `json:"video_id"  gorm:"column:video_id"`
	Resolution int    `json:"resolution"  gorm:"column:resolution"`
	Dst        string `json:"dst" gorm:"column:dst;type:varchar(255)"`
	Status     int    `json:"method"  gorm:"column:status"`
	Type       string `json:"path" gorm:"column:type;type:varchar(255)"`
}

func (TranscodingTask) TableName() string {
	return "lv_transcoding_task"
}
