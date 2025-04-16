package sundry

import (
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
)

// Upload 存储某一类型保存本地时的保存目录
type Upload struct {
	gorm.Model
	//上传文件的类型interface:
	//			videoContributionCover;videoContribution;
	//			articleContribution;articleContributionCover;
	//			createFavoritesCover;userAvatar;liveCover
	Interfaces string  `json:"interface"  gorm:"column:interface"`
	Method     string  `json:"method"  gorm:"column:method"`   // 上传文件方式local、aliyunOss
	Path       string  `json:"path" gorm:"column:path"`        // 要保存视频文件的目录路径
	Quality    float64 `json:"quality"  gorm:"column:quality"` // 图片文件的保存大小
}

func (Upload) TableName() string {
	return "lv_upload_method"
}

// IsExistByField 查找某个字段值的记录
func (upd *Upload) IsExistByField(field string, value any) bool {
	err := global.MysqlDb.Model(&Upload{}).Where(field+" = ?", value).First(&upd).Error
	return err == nil
}
