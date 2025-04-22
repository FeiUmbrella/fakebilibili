package video

import (
	"fakebilibili/infrastructure/model/user/notice"
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
)

type Likes struct {
	gorm.Model
	Uid     uint `json:"uid" gorm:"column:uid"`
	VideoID uint `json:"video_id"  gorm:"column:video_id"`
}

type LikesList []Likes

func (Likes) TableName() string {
	return "lv_video_contribution_like"
}

// IsLike 判断用户是否点赞视频
func (lk *Likes) IsLike(uid uint, vID uint) bool {
	err := global.MysqlDb.
		Where("uid = ? AND video_id = ?", uid, vID).
		First(&lk).Error
	return err == nil
}

// Like uid给属于videoUid的视频videoID点赞
// 如果没有给视频点赞则是点赞，如果已经点过赞就是取消点赞
func (lk *Likes) Like(uid uint, videoID uint, videoUid uint) error {
	err := global.MysqlDb.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("uid = ? AND video_id = ?", uid, videoID).Find(lk).Error
		if err != nil {
			return err
		}
		if lk.ID > 0 { // 已经存在uid给videoID点过赞，表明这次行为是取消点赞
			err = tx.Where("uid = ? AND video_id = ?", uid, videoID).Delete(lk).Error
			if err != nil {
				return err
			}
			// 删除点赞信息的通知
			if videoUid == uid {
				// 给自己点赞时没有进行通知
				return nil
			}

			nt := new(notice.Notice)
			err = nt.Delete(videoUid, uid, videoID, notice.VideoLike)
			if err != nil {
				return err
			}
		} else {
			// 这次行为是点赞
			lk.Uid = uid
			lk.VideoID = videoID
			err = global.MysqlDb.Create(&lk).Error
			if err != nil {
				return err
			}
			// 点赞自己作品不进行通知
			if videoUid == uid {
				return nil
			}

			// 添加消息通知视频作者
			nt := new(notice.Notice)
			err = nt.AddNotice(videoUid, uid, videoID, notice.VideoLike, "攒了您的作品")
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
