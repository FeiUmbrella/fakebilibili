package video

import (
	"encoding/json"
	"errors"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/common"
	"fakebilibili/infrastructure/model/contribution/video/barrage"
	"fakebilibili/infrastructure/model/contribution/video/comments"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

type VideosContribution struct {
	gorm.Model
	Uid           uint           `json:"uid" gorm:"column:uid"`
	Title         string         `json:"title" gorm:"column:title;type:varchar(255)"`
	Video         datatypes.JSON `json:"video" gorm:"column:video"`           // 视频的存储路径
	Video720p     datatypes.JSON `json:"video_720p" gorm:"column:video_720p"` // {type-local/oss; path-720p的存储路径}
	Video480p     datatypes.JSON `json:"video_480p" gorm:"column:video_480p"`
	Video360p     datatypes.JSON `json:"video_360p" gorm:"column:video_360p"`
	MediaID       string         `json:"media_id" gorm:"column:media_id; type:varchar(255)"`   // oss视频文件注册媒资返回的ID
	Cover         datatypes.JSON `json:"cover" gorm:"column:cover"`                            // 封面
	VideoDuration int64          `json:"video_duration" gorm:"column:video_duration"`          // 时长
	Reprinted     int8           `json:"reprinted" gorm:"column:reprinted"`                    // todo:这个字段是干嘛的？
	Label         string         `json:"label" gorm:"column:label; type:varchar(255)"`         // 视频标签，可用于分类和查找
	Introduce     string         `json:"introduce" gorm:"column:introduce; type:varchar(255)"` // 视频介绍
	Heat          int            `json:"heat" gorm:"column:heat"`                              // 视频热度(播放量)
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
		for i, j := 0, 5-len(appendRandomVideos); i < j; i++ {
			//fmt.Println(i)
			appendRandomVideos = append(appendRandomVideos, skipVideos[rand.Intn(len(skipVideos))])
		}
	}

	*vcl = append(orderVideos, appendRandomVideos...)
	return nil
}

// GetRecommendList 获取给用户推荐的视频
func (vcl *VideosContributionList) GetRecommendList(uid uint) error {
	// 将每个用户的推荐视频放在redis中缓存30s
	res, err := global.RedisDb.Get(fmt.Sprintf("%s_%s", consts.RecommendVideosList, strconv.FormatUint(uint64(uid), 10))).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		//redis出错了，日志报告一下，然后继续查数据库
		global.Logger.Errorf("获取用户推荐视频时，查询redis出错：%v", err)
	}
	if len(res) != 0 {
		err := json.Unmarshal([]byte(res), vcl) // 将redis中缓存的推荐视频列表反序列化到vcl中
		if err != nil {
			global.Logger.Errorf("类型转换出错:%v", err)
		}
		global.Logger.Infof("请求推荐视频数据，使用了redis缓存")
		return nil
	}

	// 从数据库中取出7条视频
	err = global.MysqlDb.Preload("Likes").
		Preload("Comments").
		Preload("Barrage").
		Preload("UserInfo").
		Order("created_at desc").
		Limit(7).Find(&vcl).Error
	if err != nil {
		global.Logger.Errorf("查询数据库推荐视频出错：%v", err)
	}
	data, _ := json.Marshal(vcl)
	global.RedisDb.Set(fmt.Sprintf("%s_%s", consts.RecommendVideosList, strconv.FormatUint(uint64(uid), 10)), string(data), 30*time.Second)
	return nil
}

// Search 获取视频title中包含keyword的视频
func (vcl *VideosContributionList) Search(info common.PageInfo) error {
	return global.MysqlDb.Where("`title` LIKE ?", "%"+info.Keyword+"%").
		Preload("Likes").
		Preload("Comments").
		Preload("Barrage").
		Preload("UserInfo").
		Limit(info.Size).
		Offset((info.Page - 1) * info.Size).
		Order("created_at desc").
		Find(&vcl).Error

}

// GetVideoComments 获取视频评论
func (vc *VideosContribution) GetVideoComments(vId uint, pageInfo common.PageInfo) bool {
	/*
		这里在预加载Comments的时候，传入了函数来达到更精细的控制，传入函数的作用为：
		预加载Comments，预加载与Comments相关的UserInfo
		将Comments排序并分页
	*/
	err := global.MysqlDb.Where("id = ?", vId).
		Preload("Likes").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserInfo").
				Order("created_at desc").
				Limit(pageInfo.Size).
				Offset((pageInfo.Page - 1) * pageInfo.Size)
		}).Find(&vc).Error
	return err == nil
}

// FindByID 获取视频
func (vc *VideosContribution) FindByID(vid uint) error {
	return global.MysqlDb.Where("id = ?", vid).
		Preload("Likes").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserInfo").
				Order("created_at desc")
		}).Preload("Barrage").
		Preload("UserInfo").
		Order("created_at desc").
		Find(&vc).Error
}

// Watch 添加播放量
func (vc *VideosContribution) Watch(vid uint) error {
	return global.MysqlDb.Model(&VideosContribution{}).
		Where("id = ?", vid).
		Updates(map[string]interface{}{
			"heat": gorm.Expr("Heat + ?", 1),
		}).Error
}

// Create 在数据库中创建视频信息
func (vc *VideosContribution) Create() bool {
	err := global.MysqlDb.Create(&vc).Error
	return err == nil
}

// Save 保存视频信息
// save: 未传入主键则进行Insert操作，存在主键进行update操作
func (vc *VideosContribution) Save() bool {
	err := global.MysqlDb.Save(vc).Error
	return err == nil
}

// Update 更新视频信息
func (vc *VideosContribution) Update(fields map[string]interface{}) bool {
	err := global.MysqlDb.Model(vc).Updates(fields).Error
	return err == nil
}

// Delete 删除视频信息
func (vc *VideosContribution) Delete(id uint, uid uint) bool {
	err := global.MysqlDb.Where("id = ?", id).Find(&vc).Error
	if err != nil {
		return false
	}
	if vc.Uid != uid {
		return false
	}
	err = global.MysqlDb.Delete(&vc).Error
	return err == nil
}

// GetVideoManagementList 获取个人发布的视频和评论信息
func (vcl *VideosContributionList) GetVideoManagementList(pageInfo common.PageInfo, uid uint) error {
	return global.MysqlDb.Where("uid = ?", uid).
		Preload("Likes").
		Preload("Comments").
		Preload("Barrage").
		Limit(pageInfo.Size).
		Offset((pageInfo.Page - 1) * pageInfo.Size).
		Order("created_at desc").
		Find(&vcl).Error
}
