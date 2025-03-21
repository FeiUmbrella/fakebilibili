package video

type WatchRecord struct {
	Id           int64  `json:"id" gorm:"column:id"`
	Uid          uint   `json:"uid" gorm:"column:uid"`
	VideoID      uint   `json:"video_id"  gorm:"column:video_id"`
	WatchTime    string `json:"watch_time" gorm:"watch_time;type:varchar(255)"`
	CreateTime   string `json:"create_time" gorm:"create_time;type:varchar(255)"`
	DeleteStatus int    `json:"delete_status" gorm:"delete_status"`
}

func (WatchRecord) TableName() string {
	return "lv_watch_record"
}
