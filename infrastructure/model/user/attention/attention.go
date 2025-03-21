package attention

import (
	user2 "fakebilibili/infrastructure/model/user"
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
