package video

import (
	"errors"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/common"
	"fakebilibili/infrastructure/model/contribution/video/barrage"
	"fakebilibili/infrastructure/model/contribution/video/comments"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"math/rand"
)

type VideosContribution struct {
	gorm.Model
	Uid           uint           `json:"uid" gorm:"column:uid"`
	Title         string         `json:"title" gorm:"column:title;type:varchar(255)"`
	Video         datatypes.JSON `json:"video" gorm:"column:video"` //默认1080p
	Video720p     datatypes.JSON `json:"video_720p" gorm:"column:video_720p"`
	Video480p     datatypes.JSON `json:"video_480p" gorm:"column:video_480p"`
	Video360p     datatypes.JSON `json:"video_360p" gorm:"column:video_360p"`
	MediaID       string         `json:"media_id" gorm:"column:media_id; type:varchar(255)"`
	Cover         datatypes.JSON `json:"cover" gorm:"column:cover"`
	VideoDuration int64          `json:"video_duration" gorm:"column:video_duration"`
	Reprinted     int8           `json:"reprinted" gorm:"column:reprinted"`
	Label         string         `json:"label" gorm:"column:label; type:varchar(255)"`
	Introduce     string         `json:"introduce" gorm:"column:introduce; type:varchar(255)"`
	Heat          int            `json:"heat" gorm:"column:heat"`
	//todo:加了一个visible字段，可能会引起很多连锁反应
	IsVisible int `json:"is_visible" gorm:"column:is_visible"`

	UserInfo user.User            `json:"user_info" gorm:"foreignKey:Uid"`
	Likes    LikesList            `json:"likes" gorm:"foreignKey:VideoID" `
	Comments comments.CommentList `json:"comments" gorm:"foreignKey:VideoID"`
	Barrage  barrage.BarragesList `json:"barrage" gorm:"foreignKey:VideoID"`
}

type VideosContributionList []VideosContribution

func (VideosContribution) TableName() string {
	return "lv_video_contribution"
}

// GetVideoListBySpace 获取空间视频列表
func (vcl *VideosContributionList) GetVideoListBySpace(id uint) error {
	return global.MysqlDb.
		Model(&VideosContribution{}).
		Where("uid = ?", id).
		Preload("Likes").
		Preload("Comments").
		Preload("Barrage").
		Order("created_at desc").
		Find(&vcl).Error
}

// GetHomeVideoList 获取主页推荐视频
func (vcl *VideosContributionList) GetHomeVideoList(pageInfo common.PageInfo) error {
	var offset int
	// todo:这里的按照页数来加载首页视频个数没看懂
	if pageInfo.Page == 1 {
		pageInfo.Size = 11
		offset = (pageInfo.Page - 1) * pageInfo.Size
	}
	offset = (pageInfo.Page-2)*pageInfo.Size + 11

	// 按照热度排序查询出pageInfo.size-5条视频，再随机拼凑5条视频
	var orderVideos []VideosContribution
	orderSize := pageInfo.Size - 5
	if orderSize > 0 {
		err := global.MysqlDb.Preload("Likes").
			Preload("Comments").
			Preload("Barrage").
			Preload("UserInfo").
			Where("is_visible = ?", 1).
			Limit(pageInfo.Size).Offset(offset).
			Order("heat desc").
			Find(&orderVideos).Error
		if err != nil {
			return errors.New("failed to query videos  order by heat desc:" + err.Error())
		}
	}
	// 将按照热度选出的视频id保存在 Redis的bitmap中
	for _, video := range orderVideos {
		_, err := global.RedisDb.SetBit(fmt.Sprintf("%s%d", consts.UniqueVideoRecommendPrefix, -1), int64(video.ID), 1).Result()
		if err != nil {
			global.Logger.Errorf("set bitmap of orderVideos at home page  failed:%v", err)
		}
	}
	// 随机再取10条视频，找到上面未取到的5条视频
	var randomVideos []VideosContribution
	if err := global.MysqlDb.Preload("Comments").
		Preload("Likes").
		Preload("Barrage").
		Preload("UserInfo").
		Where("is_visible=?", 1).
		Order("RAND()").Limit(10).
		Find(&randomVideos).Error; err != nil {
		return errors.New("failed to query videos  randomly:" + err.Error())
	}

	// 去重
	var appendRandomVideos []VideosContribution
	var skipVideos []VideosContribution
	for _, video := range randomVideos {
		result, err := global.RedisDb.GetBit(fmt.Sprintf("%s%d", consts.UniqueVideoRecommendPrefix, -1), int64(video.ID)).Result()
		if err != nil {
			global.Logger.Errorf("get bitmap failed:%v", err)
		}
		if result == 1 { // 视频已经在选中过
			skipVideos = append(skipVideos, video)
		} else {
			appendRandomVideos = append(appendRandomVideos, video)
		}
	}

	// 去重后不足5条视频,从筛掉的视频中随机填充到5个
	if len(appendRandomVideos) < 5 {
		for i := 0; i < 5-len(appendRandomVideos); i++ {
			appendRandomVideos = append(appendRandomVideos, skipVideos[rand.Intn(len(skipVideos))])
		}
	}

	*vcl = append(orderVideos, appendRandomVideos...)
	return nil
}
