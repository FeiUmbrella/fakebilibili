package article

import "gorm.io/gorm"

// Likes 对文章的赞
type Likes struct {
	gorm.Model
	Uid       uint `json:"uid"`
	ArticleID uint `json:"article_id" gorm:"article_id"`
}

type LikesList []Likes

func (Likes) TableName() string {
	return "lv_article_contribution_likes"
}
