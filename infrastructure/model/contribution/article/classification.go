package article

import (
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
)

type Classification struct {
	gorm.Model
	AID   uint   `json:"a_id" gorm:"column:a_id"`
	Label string `json:"label" gorm:"type:varchar(255)"`
}

type ClassificationsList []Classification

func (Classification) TableName() string {
	return "lv_article_classification"
}

// FindAll 返回所有文章
func (cl *ClassificationsList) FindAll() error {
	return global.MysqlDb.Find(cl).Error
}
