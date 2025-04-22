package comments

import (
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/notice"
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Uid            uint   `json:"uid"`
	VideoID        uint   `json:"video_id" gorm:"video_id"`
	Context        string `json:"context" gorm:"type:text"`
	CommentID      uint   `json:"comment_id" gorm:"comment_id"` // (父评论)该条评论是在CommentID评论下的 CommentID=0为根评论
	CommentUserID  uint   `json:"comment_user_id" gorm:"comment_user_id"`
	CommentFirstID uint   `json:"comment_first_id" gorm:"comment_first_id"` // 这条评论的根评论
	Heat           int    `json:"heat" gorm:"heat"`                         // 被赞数即热度

	UserInfo  user.User `json:"user_info" gorm:"foreignKey:Uid"`
	VideoInfo VideoInfo `json:"video_info" gorm:"foreignKey:VideoID"`
}
type CommentList []Comment

func (Comment) TableName() string {
	return "lv_video_contribution_comments"
}

// VideoInfo 临时加一个video模型解决依赖循环
type VideoInfo struct {
	gorm.Model
	Uid   uint           `json:"uid" gorm:"uid"`
	Title string         `json:"title" gorm:"title"`
	Video datatypes.JSON `json:"video" gorm:"video"`
	Cover datatypes.JSON `json:"cover" gorm:"cover"`
}

func (VideoInfo) TableName() string {
	return "lv_video_contribution"
}

// Find 更加id查询某一条评论
func (c *Comment) Find(id uint) {
	_ = global.MysqlDb.Where("id = ?", id).First(&c).Error
}

// GetCommentFirstID 获取子评论所属的根评论
func (c *Comment) GetCommentFirstID() uint {
	_ = global.MysqlDb.Where("id = ?", c.ID).Find(&c).Error
	if c.CommentID != 0 {
		c.ID = c.CommentID
		c.GetCommentFirstID()
	}
	return c.ID
}

// GetCommentUserID 找到评论所属用户
func (c *Comment) GetCommentUserID() uint {
	_ = global.MysqlDb.Where("id = ?", c.ID).Find(&c).Error
	return c.Uid
}

// Create 创建评论
func (c *Comment) Create() bool {
	err := global.MysqlDb.Transaction(func(tx *gorm.DB) error {
		videoInfo := new(VideoInfo)
		err := tx.Where("id = ?", c.VideoID).Find(videoInfo).Error
		if err != nil {
			return err
		}
		err = tx.Create(&c).Error
		if err != nil {
			return err
		}
		//消息通知
		if videoInfo.Uid == c.Uid {
			// 给自己视频评论，不进行notice通知
			return nil
		}
		// 添加消息通知
		ne := new(notice.Notice)
		err = ne.AddNotice(videoInfo.Uid, c.Uid, videoInfo.ID, notice.VideoComment, c.Context)
		if err != nil {
			return err
		}
		return nil
	})
	return err == nil
}
