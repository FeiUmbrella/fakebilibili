package chat

import (
	"errors"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
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

// GetListByID 查找用户的聊天列表
func (cl *ChatList) GetListByID(uid uint) error {
	return global.MysqlDb.
		Where("uid = ?", uid).
		Preload("ToUserInfo").
		Order("updated_at desc").
		Find(cl).Error
}

// AddChat 创建一条列表
func (cl *ChatsListInfo) AddChat() error {
	// 判断聊天列表是否存在两人
	temp := new(ChatsListInfo)
	err := global.MysqlDb.Where("uid = ? AND tid = ?", cl.Uid, cl.Tid).First(temp).Error
	if err == nil { // 存在,就更新一下最新消息和时间
		global.MysqlDb.Model(&ChatsListInfo{}).
			Updates(map[string]interface{}{
				"last_at":      cl.LastAt,
				"last_message": cl.LastMessage,
			})
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) { // 不存在则创建
		return global.MysqlDb.Create(cl).Error
	} else { // 其他错误
		return err
	}
}

// DeleteChat 删除聊天记录
func (cl *ChatsListInfo) DeleteChat(tid, uid uint) error {
	return global.MysqlDb.Where("uid = ? AND tid = ?", uid, tid).Delete(cl).Error
}

// GetUnreadNumber 获取用户未读私信数量
func (cl *ChatsListInfo) GetUnreadNumber(uid uint) *int64 {
	var num int64
	err := global.MysqlDb.Model(&ChatsListInfo{}).
		Select("SUM(IFNULL(unread,0)) as total_unread").
		Where("uid = ?", uid).Scan(&num).Error
	if err != nil {
		fmt.Println(err)
	}
	return &num
}
