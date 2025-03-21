package chat

import (
	"fakebilibili/infrastructure/model/user"
	"gorm.io/gorm"
)

// Msg 两个用户发送的信息
type Msg struct {
	gorm.Model
	Uid     uint   `json:"uid"`
	Tid     uint   `json:"tid"`
	Type    string `json:"type" gorm:"type:varchar(255)"`
	Message string `json:"message" gorm:"type:text"`

	UInfo user.User `json:"UInfo" gorm:"foreignKey:uid"`
	TInfo user.User `json:"TInfo" gorm:"foreignKey:tid"`
}

type MsgList []Msg

func (Msg) TableName() string {
	return "lv_users_chat_msg"
}
