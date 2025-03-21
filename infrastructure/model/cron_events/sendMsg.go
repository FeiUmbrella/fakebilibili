package cron_events

// todo:这个表示干啥使得？
type CronEvent struct {
	Id       int64  `json:"id" gorm:"id"`
	LastTime string `json:"last_time" gorm:"last_time;type:varchar(255)"` //上一次扫表的时间
}

func (CronEvent) TableName() string {
	return "lv_cron_events"
}
