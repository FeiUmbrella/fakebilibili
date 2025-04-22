package collect

import (
	"fakebilibili/infrastructure/model/contribution/video"
	user2 "fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
)

// Collect 某个收藏夹与被收藏视频、已经所属用户的关联信息
type Collect struct {
	gorm.Model
	Uid         uint `json:"uid"`
	FavoritesID uint `json:"favorites_id" gorm:"column:favorites_id"`
	VideoID     uint `json:"video_id" gorm:"column:video_id"`

	UserInfo  user2.User               `json:"userInfo" gorm:"foreignKey:Uid"`
	VideoInfo video.VideosContribution `json:"videoInfo" gorm:"foreignKey:VideoID"`
}

type CollectsList []Collect

func (Collect) TableName() string {
	return "lv_users_collect"
}

// DeleteByFavoritesID 删除某一收藏夹中的所有收藏内容
func (c *Collect) DeleteByFavoritesID(id uint) bool {
	err := global.MysqlDb.Model(&Collect{}).Where("favorites_id=?", id).Delete(&Collect{}).Error
	return err == nil
}

// FindFavoriteIncludeVideo 找到包含某个视频的所有收藏夹id
func (cl *CollectsList) FindFavoriteIncludeVideo(vid uint) error {
	return global.MysqlDb.Model(&Collect{}).
		Where("video_id=?", vid).
		Find(&cl).Error
}

// DeleteOneVideoInFavorite 删除在一个收藏夹中的一个视频
func (c *Collect) DeleteOneVideoInFavorite(vid, fid uint) error {
	return global.MysqlDb.Model(&Collect{}).
		Where("video_id=? AND favorites_id=?", vid, fid).
		Delete(&Collect{}).Error
}

// FindOneVideoInFavorite 查找在一个收藏夹中的一个视频的记录
func (c *Collect) FindOneVideoInFavorite(vid, fid uint) error {
	err := global.MysqlDb.Model(&Collect{}).
		Where("video_id=? AND favorites_id=?", vid, fid).
		First(&c).Error
	//fmt.Printf("错误类型: %T, 错误详情: %v\n", err, err)
	return err
}

// Create 创建一条收藏视频的记录
func (c *Collect) Create() bool {
	err := global.MysqlDb.Create(&c).Error
	return err == nil
}

// FindVideosByFavoriteID 查找一个收藏夹中收藏的视频
func (cl *CollectsList) FindVideosByFavoriteID(fid uint) error {
	return global.MysqlDb.
		Where("favorites_id=?", fid).
		Preload("VideoInfo").
		Find(&cl).Error
}

// FindIsCollectByFavorites 判断某一用户的所有收藏夹中是否收藏了视频
func (cl *CollectsList) FindIsCollectByFavorites(vID uint, fIDs []uint) bool {
	// 没创建收藏夹，直接return false
	if len(fIDs) == 0 {
		return false
	}
	err := global.MysqlDb.
		Where("video_id=?", vID).
		Where("favorites_id IN (?)", fIDs).
		Find(&cl).Error
	if err != nil || len(*cl) == 0 {
		return false
	}
	return true
}
