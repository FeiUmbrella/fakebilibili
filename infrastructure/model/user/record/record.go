package record

import (
	"errors"
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

	/*
		todo：
			这里外键设置有问题，与VideoInfo和ArticleInfo的外键都是to_id，
			导致一条(uid=2, type=video, to_id=1)的record，也同时指向article_id=1的文章
			如果要插入上面这条record，如果id=1的视频存在但id=1的文章不存在就会导致插入失败！
			可以将外键分开，搞一个to_video_id，一个to_article_id分别作为VideoInfo和Article的外键
	*/
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

// AddLiveRecord 添加观看直播的记录 uid是用户，roomId是直播间ID也是直播间主播ID
func (r *Record) AddLiveRecord(uid uint, roomId uint) error {
	err := global.MysqlDb.Model(&Record{}).
		Where("uid = ? AND to_id = ? AND type = ?", uid, roomId, "live").
		First(r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建记录
		r.Uid = uid
		r.Type = "live"
		r.ToId = roomId
		return global.MysqlDb.Create(r).Error
	}
	if err != nil {
		return err
	}
	// 存在记录，更新一下
	return global.MysqlDb.Where("id = ?", r.ID).Updates(r).Error
}

// AddVideoRecord 添加观看视频的记录
func (r *Record) AddVideoRecord(uid uint, videoId uint) error {
	err := global.MysqlDb.Model(&Record{}).
		Where("uid = ? AND to_id = ? AND type = ?", uid, videoId, "video").
		First(r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建记录
		r.Uid = uid
		r.Type = "video"
		r.ToId = videoId
		return global.MysqlDb.Create(r).Error
	}
	if err != nil {
		return err
	}
	// 存在记录，更新一下
	return global.MysqlDb.Where("id = ?", r.ID).Updates(r).Error
}

// AddArticleRecord 添加观看文章的记录
func (r *Record) AddArticleRecord(uid uint, articleId uint) error {
	err := global.MysqlDb.Model(&Record{}).
		Where("uid = ? AND to_id = ? AND type = ?", uid, articleId, "article").
		First(r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建记录
		r.Uid = uid
		r.Type = "article"
		r.ToId = articleId
		return global.MysqlDb.Create(r).Error
	}
	if err != nil {
		return err
	}
	// 存在记录，更新一下
	return global.MysqlDb.Where("id = ?", r.ID).Updates(r).Error
}
