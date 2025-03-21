package chat

import (
	"fakebilibili/infrastructure/model/user"
	"gorm.io/gorm"
	"time"
)

// ChatsListInfo 聊天列表
type ChatsListInfo struct {
	gorm.Model
	Uid         uint      `json:"uid"`                              // 用户
	Tid         uint      `json:"tid"`                              // 聊天对象
	Unread      int       `json:"unread"`                           // 是否已读
	LastMessage string    `json:"last_message" gorm:"last_message"` // todo: 最后一条信息？
	LastAt      time.Time `json:"last_at" gorm:"last_at"`           // todo：最后一条信息时间？

	ToUserInfo user.User `json:"toUserInfo" gorm:"foreignKey:tid"` // 链接聊条对象用户信息
}

type ChatList []ChatsListInfo

func (ChatsListInfo) TableName() string {
	return "lv_users_chat_list"
}
