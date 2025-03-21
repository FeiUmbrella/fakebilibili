package article

import (
	"fakebilibili/infrastructure/model/user"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ArticlesContribution 贡献的文章
type ArticlesContribution struct {
	gorm.Model
	Uid                uint           `json:"uid"`
	ClassificationID   uint           `json:"classification_id" gorm:"classification_id"`
	Title              string         `json:"title" gorm:"type:varchar(255)"`
	Cover              datatypes.JSON `json:"cover"`
	Label              string         `json:"label" gorm:"type:varchar(255)"`
	Content            string         `json:"content" gorm:"type:text"`
	ContentStorageType string         `json:"content_storage_type" gorm:"content_storage_type;type:varchar(255)"`
	IsComments         int8           `json:"is_comments" gorm:"is_comments"` // todo: 这个字段干啥？
	Heat               int            `json:"heat" gorm:"heat"`

	// 外键关联表
	UserInfo       user.User      `json:"user_info" gorm:"foreignKey:Uid"`
	Likes          LikesList      `json:"likes" gorm:"foreignKey:ArticleID"`
	Comments       CommentList    `json:"comments" gorm:"foreignKey:ArticleID"`
	Classification Classification `json:"classification" gorm:"foreignKey:ClassificationID"`
}

type ArticlesContributionList []ArticlesContribution

func (ArticlesContribution) TableName() string {
	return "lv_article_contribution"
}
