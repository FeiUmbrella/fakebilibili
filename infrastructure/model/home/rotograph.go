package home

import (
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Rotograph 主页的视频轮播图
type Rotograph struct {
	gorm.Model
	Title string         `json:"title" gorm:"column:title;type:varchar(255)"`
	Cover datatypes.JSON `json:"cover" gorm:"column:cover"`
	Color string         `json:"color" gorm:"column:color;type:varchar(255)" `
	Type  string         `json:"type" gorm:"column:type;type:varchar(255)"`
	ToId  uint           `json:"to_id" gorm:"column:to_id"` // todo:这个什么意思？
}

type List []Rotograph

func (Rotograph) TableName() string {
	return "lv_home_rotograph"
}

// GetALL 获取轮播图
func (l *List) GetALL() error {
	return global.MysqlDb.Find(&l).Error
}

// Create 创建轮播图记录
func (r *Rotograph) Create(tx *gorm.DB) error {
	return tx.Model(&Rotograph{}).Create(r).Error
}
