package article

import (
	"fakebilibili/infrastructure/model/common"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
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

// Create 创建文章信息
func (ac *ArticlesContribution) Create() bool {
	err := global.MysqlDb.Create(&ac).Error
	return err == nil
}

// Delete 删除文章信息
func (ac *ArticlesContribution) Delete(aid, uid uint) bool {
	err := global.MysqlDb.Where("id = ?", aid).Find(ac).Error
	if err != nil {
		return false
	}
	if ac.Uid != uid {
		return false
	}
	err = global.MysqlDb.Delete(ac).Error
	return err == nil
}

// Update 更新文章信息
func (ac *ArticlesContribution) Update(info map[string]interface{}) bool {
	err := global.MysqlDb.Model(ac).Where("id = ?", ac.ID).Updates(info).Error
	return err == nil
}

// GetArticleBySpace 获取空间专栏
func (acl *ArticlesContributionList) GetArticleBySpace(id uint) error {
	return global.MysqlDb.Where("uid = ?", id).
		Preload("Likes").
		Preload("Comments").
		Preload("Classification").
		Order("created_at desc").
		Find(&acl).Error
}

// GetList 获取文章
func (acl *ArticlesContributionList) GetList(pageInfo common.PageInfo) bool {
	err := global.MysqlDb.Preload("Likes").
		Preload("Classification").
		Preload("UserInfo").
		Preload("Comments").
		Limit(pageInfo.Size).Offset((pageInfo.Page - 1) * pageInfo.Size).
		Order("created_at desc").
		Find(acl).Error
	return err == nil
}

// GetListByUid 获取用户文章
func (acl *ArticlesContributionList) GetListByUid(uid uint) bool {
	err := global.MysqlDb.
		Where("uid = ?", uid).
		Preload("Likes").
		Preload("Classification").
		Preload("UserInfo").
		Preload("Comments").
		Order("created_at desc").
		Find(acl).Error
	return err == nil
}

// GetArticleComments 获取文章评论
func (ac *ArticlesContribution) GetArticleComments(aid uint, pageInfo common.PageInfo) bool {
	err := global.MysqlDb.
		Where("id = ?", aid).
		Preload("Likes").
		Preload("Classification").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserInfo").
				Limit(pageInfo.Size).Offset((pageInfo.Page - 1) * pageInfo.Size).
				Order("created_at desc")
		}).Find(ac).Error
	return err == nil
}

// GetAllCount 获取所有文章数量
func (acl *ArticlesContributionList) GetAllCount(cnt *int64) bool {
	err := global.MysqlDb.Find(acl).Count(cnt).Error
	return err == nil
}

// GetInfoByID 查询单个文章
func (ac *ArticlesContribution) GetInfoByID(aid uint) bool {
	err := global.MysqlDb.
		Where("id = ?", aid).
		Preload("Likes").
		Preload("Classification").
		Preload("UserInfo").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserInfo").
				Order("created_at desc")
		}).Find(ac).Error
	return err == nil
}

// Watch 添加观看次数
func (ac *ArticlesContribution) Watch(id uint) error {
	return global.MysqlDb.Model(ac).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"heat": gorm.Expr("Heat + ?", 1),
		}).Error
}

// GetArticleManagementList 创作空间获取个人发布专栏
func (acl *ArticlesContributionList) GetArticleManagementList(pageInfo common.PageInfo, uid uint) error {
	err := global.MysqlDb.Where("uid = ?", uid).
		Preload("Likes").
		Preload("Classification").
		Preload("Comments").
		Limit(pageInfo.Size).Offset((pageInfo.Page - 1) * pageInfo.Size).
		Order("created_at desc").
		Find(acl).Error
	return err
}
