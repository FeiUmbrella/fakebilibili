package attention

import (
	user2 "fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
)

// Attention 关注
type Attention struct {
	gorm.Model
	Uid         uint `json:"uid" gorm:"column:uid"`                   // 用户A
	AttentionID uint `json:"attention_id" gorm:"column:attention_id"` // A 关注的人B

	UserInfo          user2.User `json:"user_info" gorm:"foreignKey:Uid"`                   // 用uid作为外键，关联A的信息
	AttentionUserInfo user2.User `json:"attention_user_info" gorm:"foreignKey:AttentionID"` // 关联B的信息
}

type AttentionsList []Attention

func (Attention) TableName() string {
	return "lv_users_attention"
}

// GetAttentionList 获取uid的关注列表
func (atl *AttentionsList) GetAttentionList(uid uint) error {
	err := global.MysqlDb.Model(&Attention{}).Preload("AttentionUserInfo").Where("uid = ?", uid).Find(&atl).Error
	return err
}

// GetVermicelliList 获取uid的粉丝列表
func (atl *AttentionsList) GetVermicelliList(uid uint) error {
	err := global.MysqlDb.Model(&Attention{}).Preload("UserInfo").Where("attention_id = ?", uid).Find(&atl).Error
	return err
}

// GetAttentionListByIdArr 获取uid关注列表中每个关注者的id
func (atl *AttentionsList) GetAttentionListByIdArr(uid uint) ([]uint, error) {
	arr := make([]uint, 0)
	// 找到uid关注的所有人
	err := global.MysqlDb.Model(&Attention{}).Where("uid = ?", uid).Find(&atl).Error
	if err != nil {
		return arr, err
	}
	// 提取这些人的id
	for _, at := range *atl {
		arr = append(arr, at.AttentionID)
	}
	return arr, nil
}

// IsAttention 查找uid是否关注aid
func (at *Attention) IsAttention(uid uint, aid uint) bool {
	err := global.MysqlDb.Model(&Attention{}).Where("uid = ? AND attention_id = ?", uid, aid).First(&at).Error
	return err == nil
}

// GetAttentionNum 得到关注数量
func (at *Attention) GetAttentionNum(uid uint) (int64, error) {
	var cnt int64
	err := global.MysqlDb.Model(&Attention{}).Where("uid = ?", uid).Count(&cnt).Error
	return cnt, err
}

// GetVermicelliNum 得到粉丝数量
func (at *Attention) GetVermicelliNum(uid uint) (int64, error) {
	var cnt int64
	err := global.MysqlDb.Model(&Attention{}).Where("attention_id = ?", uid).Count(&cnt).Error
	return cnt, err
}
