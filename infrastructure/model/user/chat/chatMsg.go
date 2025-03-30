package chat

import (
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"gorm.io/gorm"
	"time"
)

// Msg 两个用户发送的信息
type Msg struct {
	gorm.Model
	Uid     uint   `json:"uid"`
	Tid     uint   `json:"tid"`
	Type    string `json:"type" gorm:"type:varchar(255)"` // todo:两个人的聊题记录有什么类型？
	Message string `json:"message" gorm:"type:text"`

	UInfo user.User `json:"UInfo" gorm:"foreignKey:uid"`
	TInfo user.User `json:"TInfo" gorm:"foreignKey:tid"`
}

type MsgList []Msg

func (Msg) TableName() string {
	return "lv_users_chat_msg"
}

// FindList 查找两人的聊天信息
func (mgl *MsgList) FindList(uid, tid uint) error {
	return global.MysqlDb.
		Where("(uid = ? AND tid = ?) OR (uid = ? AND tid = ?)", uid, tid, tid, uid).
		Preload("UInfo").
		Preload("TInfo").
		Order("created_at desc").
		Limit(30).
		Find(mgl).Error
}

// FindHistoryMst 查询某天之前的历史聊天记录
func (mgl *MsgList) FindHistoryMst(uid, tid uint, lastTime time.Time) error {
	return global.MysqlDb.
		Debug().
		Where("(uid = ? AND tid = ?) OR (uid = ? AND tid = ?)", uid, tid, tid, uid).
		Where("created_at < ?", lastTime).
		Preload("UInfo").
		Preload("TInfo").
		Order("created_at desc").
		Limit(30).
		Find(mgl).Error
}

// GetLastMessage 查找两人最后一条聊天记录
func (mg *Msg) GetLastMessage(uid, tid uint) error {
	return global.MysqlDb.
		Where("(uid = ? AND tid = ?) OR (uid = ? AND tid = ?)", uid, tid, tid, uid).
		Order("created_at desc").
		First(mg).Error
}
