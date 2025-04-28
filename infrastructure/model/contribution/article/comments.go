package article

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
		articleInfo := new(Article)
		err := tx.Where("id = ?", c.ArticleID).Find(articleInfo).Error
		if err != nil {
			return err
		}
		err = tx.Create(&c).Error
		if err != nil {
			return err
		}
		//消息通知
		if articleInfo.Uid == c.Uid {
			// 给自己视频评论，不进行notice通知
			return nil
		}
		// 添加消息通知
		ne := new(notice.Notice)
		err = ne.AddNotice(articleInfo.Uid, c.Uid, articleInfo.ID, notice.VideoComment, c.Context)
		if err != nil {
			return err
		}
		return nil
	})
	return err == nil
}
