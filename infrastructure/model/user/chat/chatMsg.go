package chat

import (
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
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

// AddMessage 保存信息
func (mg *Msg) AddMessage() error {
	// 保存具体信息Msg，更新ChatListInfo
	// 这两个操作需要完成，使用事务完成这两个操作，如果失败事务自动回滚
	err := global.MysqlDb.Transaction(func(tx *gorm.DB) error {
		// 保存Msg
		err := tx.Create(mg).Error
		if err != nil {
			return fmt.Errorf("添加聊天记录失败")
		}

		// 聊天列表添加记录
		// uid用户和tid用户聊天列表加载的时候都是从数据库中找到对应 Uid=uid/tid 的 ChatsListInfo
		// 所以要存两条数据库记录
		// 这个是在加载uid用户聊天列表是加载的数据库记录
		uidChatListInfo := &ChatsListInfo{
			Uid:         mg.Uid,
			Tid:         mg.Tid,
			LastMessage: mg.Message,
			LastAt:      time.Now(),
		}
		err = uidChatListInfo.AddChat()
		if err != nil {
			return fmt.Errorf("添加聊天列表记录失败")
		}
		// 这个是在加载tid用户聊天列表是加载的数据库记录
		tidChatListInfo := &ChatsListInfo{
			Uid:         mg.Tid,
			Tid:         mg.Uid,
			LastMessage: mg.Message,
			LastAt:      time.Now(),
		}
		err = tidChatListInfo.AddChat()
		if err != nil {
			return fmt.Errorf("添加聊天列表记录失败")
		}

		// 成功完成事务
		return nil
	})
	// 返回事务完成情况
	return err
}

// FindByID 查找为id的记录
func (mg *Msg) FindByID(id uint) error {
	return global.MysqlDb.Model(&Msg{}).
		Where("id = ?", id).
		Preload("TInfo").
		Preload("UInfo").
		First(mg).Error
}
