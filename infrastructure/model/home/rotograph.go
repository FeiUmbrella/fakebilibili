package home

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// todo:这个表是干啥使得？
type Rotograph struct {
	gorm.Model
	Title string         `json:"title" gorm:"column:title;type:varchar(255)"`
	Cover datatypes.JSON `json:"cover" gorm:"column:cover"`
	Color string         `json:"color" gorm:"column:color;type:varchar(255)" `
	Type  string         `json:"type" gorm:"column:type;type:varchar(255)"`
	ToId  uint           `json:"to_id" gorm:"column:to_id"`
}

type List []Rotograph

func (Rotograph) TableName() string {
	return "lv_home_rotograph"
}
