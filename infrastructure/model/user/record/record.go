package record

import (
	"fakebilibili/infrastructure/model/common"
	"fakebilibili/infrastructure/model/contribution/article"
	"fakebilibili/infrastructure/model/contribution/video"
	user2 "fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
)

// Record 浏览记录
type Record struct {
	gorm.Model
	Uid  uint   `json:"uid"`
	Type string `json:"type" gorm:"type:varchar(255)"` // video/article/live
	ToId uint   `json:"to_id" gorm:"column:to_id"`     // 浏览的视频/文章 id

	VideoInfo   video.VideosContribution     `json:"videoInfo" gorm:"foreignKey:to_id"`
	UserInfo    user2.User                   `json:"userInfo" gorm:"foreignKey:uid"`
	ArticleInfo article.ArticlesContribution `json:"articleInfo" gorm:"foreignKey:to_id"`
}
type RecordList []Record

func (Record) TableName() string {
	return "lv_users_record"
}

// GetRecordListByUid 查找用户记录
func (rl *RecordList) GetRecordListByUid(uid uint, page common.PageInfo) error {
	// 查找记录的同时预加载相应的用户中的直播间信息、视频信息、文章信息
	return global.MysqlDb.Model(&Record{}).
		Where("uid = ?", uid).
		Preload("VideoInfo").
		Preload("ArticleInfo").
		Preload("UserInfo.LiveInfo").
		Limit(page.Size).
		Offset((page.Page - 1) * page.Size).
		Order("created_at desc").
		Find(rl).Error
}

// ClearRecord 清空历史记录
func (r *Record) ClearRecord(uid uint) error {
	return global.MysqlDb.
		Model(&Record{}).
		Where("uid = ?", uid).
		Delete(r).Error
}

// DeleteRecordByID 删除一条历史记录
func (r *Record) DeleteRecordByID(id, uid uint) error {
	return global.MysqlDb.
		Model(&Record{}).
		Where("id = ? AND uid = ?", id, uid).
		Delete(r).Error
}
