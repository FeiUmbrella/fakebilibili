package comments

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"os/user"
)

type Comment struct {
	gorm.Model
	Uid            uint   `json:"uid"`
	VideoID        uint   `json:"video_id" gorm:"video_id"`
	Context        string `json:"context" gorm:"type:text"`
	CommentID      uint   `json:"comment_id" gorm:"comment_id"`
	CommentUserID  uint   `json:"comment_user_id" gorm:"comment_user_id"`
	CommentFirstID uint   `json:"comment_first_id" gorm:"comment_first_id"`
	Heat           int    `json:"heat" gorm:"heat"` // 被赞数即热度

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
