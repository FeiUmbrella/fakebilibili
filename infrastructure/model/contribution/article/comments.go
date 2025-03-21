package article

import (
	"fakebilibili/infrastructure/model/user"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Uid            uint   `json:"uid"`
	ArticleID      uint   `json:"article_id" gorm:"article_id"`
	Context        string `json:"context"`
	CommentID      uint   `json:"comment_id" gorm:"comment_id"` // 在这条评论下面还有评论
	CommentUserID  uint   `json:"comment_user_id" gorm:"comment_user_id"`
	CommentFirstID uint   `json:"comment_first_id" gorm:"comment_first_id"`

	UserInfo    user.User `json:"user_info" gorm:"foreignKey:uid"`
	ArticleInfo Article   `json:"article_info" gorm:"foreignKey:article_id"`
}
type CommentList []Comment

func (Comment) TableName() string {
	return "lv_article_contribution_comments"
}

type Article struct {
	gorm.Model
	Uid              uint           `json:"uid" gorm:"uid"`
	ClassificationID uint           `json:"classification_id"  gorm:"classification_id"`
	Title            string         `json:"title" gorm:"title"`
	Cover            datatypes.JSON `json:"cover" gorm:"cover"`
}

func (Article) TableName() string {
	return "lv_article_contribution"
}
