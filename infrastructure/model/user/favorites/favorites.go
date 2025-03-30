package favorites

import (
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/collect"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Favorites 用户收藏夹信息
type Favorites struct {
	gorm.Model
	Uid     uint           `json:"uid"`                                  // 收藏夹所属用户
	Title   string         `json:"title" gorm:"type:varchar(255)"`       // 收藏夹标题
	Content string         `json:"content" gorm:"type:text"`             // 收藏夹简介
	Cover   datatypes.JSON `json:"cover" gorm:"type:json;comment:cover"` // 收藏夹封面图片链接
	Max     int            `json:"max"`                                  // 单个收藏夹最大收藏视频数

	UserInfo    user.User            `json:"userInfo" gorm:"foreignKey:Uid"`
	CollectList collect.CollectsList `json:"collectList"  gorm:"foreignKey:FavoritesID"`
}

type FavoriteList []Favorites

func (Favorites) TableName() string {
	return "lv_users_favorites"
}

// Create 创建收藏夹
func (fs *Favorites) Create() bool {
	err := global.MysqlDb.Model(&Favorites{}).Create(&fs).Error
	return err == nil
}

// Find 查找收藏夹
func (fs *Favorites) Find(id uint) bool {
	err := global.MysqlDb.Where("id = ?", id).
		Preload("CollectList").
		Order("created_at desc").
		First(&fs).Error
	return err == nil
}

// FindFavoriteByFID 由id找对应收藏夹，不预加载其他信息
func (fs *Favorites) FindFavoriteByFID(fid uint) bool {
	err := global.MysqlDb.Model(&Favorites{}).
		Where("id = ?", fid).
		First(&fs).Error
	return err == nil
}

// Update 更新收藏夹
func (fs *Favorites) Update() bool {
	err := global.MysqlDb.Model(&Favorites{}).
		Where("id = ?", fs.ID).
		Updates(&fs).Error
	return err == nil
}

// GetFavoritesList 获取用户所有收藏夹
func (fsl *FavoriteList) GetFavoritesList(uid uint) error {
	return global.MysqlDb.Model(&Favorites{}).
		Where("uid=?", uid).
		Preload("UserInfo").
		Preload("CollectList").
		Order("created_at desc").
		Find(fsl).
		Error
}

// Delete uid用户删除为id的收藏夹
func (fs *Favorites) Delete(id, uid uint) error {
	err := global.MysqlDb.Model(&Favorites{}).
		Where("id = ?", id).
		First(&fs).Error
	if err != nil {
		return fmt.Errorf("查询收藏夹失败")
	}
	if fs.Uid != uid {
		return fmt.Errorf("非收藏夹拥有者不可删除")
	}
	err = global.MysqlDb.Model(&Favorites{}).Delete(&fs).Error

	// 删除收藏夹中收藏的内容
	cl := new(collect.Collect)
	if !cl.DeleteByFavoritesID(id) {
		return fmt.Errorf("删除收藏夹内容失败")
	}
	if err != nil {
		return fmt.Errorf("删除失败")
	}
	return nil
}
