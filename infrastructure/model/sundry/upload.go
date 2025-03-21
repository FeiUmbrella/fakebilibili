package sundry

import "gorm.io/gorm"

// todo:这个表干啥使得？
type Upload struct {
	gorm.Model
	Interfaces string  `json:"interface"  gorm:"column:interface"`
	Method     string  `json:"method"  gorm:"column:method"`
	Path       string  `json:"path" gorm:"column:path"`
	Quality    float64 `json:"quality"  gorm:"column:quality"` // todo:质量表示什么？
}

func (Upload) TableName() string {
	return "lv_upload_method"
}
