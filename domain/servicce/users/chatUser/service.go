package chatUser

import (
	"fakebilibili/adapter/http/receive/socket"
	socket2 "fakebilibili/adapter/http/response/socket"
	chat2 "fakebilibili/domain/servicce/users/chat"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/user/chat"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"fakebilibili/infrastructure/pkg/utils/response"
)

// sendChatMsgText 处理前端uid用户发送的信息
func sendChatMsgText(uc *UserChannel, uid uint, tid uint, info *socket.Receive) {
	// 保存信息到本地MSg
	cm := chat.Msg{
		Uid:     uid,
		Tid:     tid,
		Type:    info.Type,
		Message: info.Data,
	}
	err := cm.AddMessage()
	if err != nil {
		// 向前端uid用户通知
		response.ErrorWs(uc.Socket, "保存信息到数据库失败")
		return
	}

	msgInfo := new(chat.Msg)
	err = msgInfo.FindByID(cm.ID)
	if err != nil {
		response.ErrorWs(uc.Socket, "发送信息失败")
		return
	}
	// uid用户的头像
	photo, _ := conversion.FormattingJsonSrc(msgInfo.UInfo.Photo)

	// 给自己发送信息不推送
	if uid == tid {
		return
	}

	if _, ok := chat2.Severe.UserMapChannel[tid]; ok {
		// tid用户在线
		if _, ok := chat2.Severe.UserMapChannel[tid].ChatList[uid]; ok {
			// tid用户正处在与uid用户的聊天窗口，直接向tid用户推送uid用户发送的信息
			response.SuccessWs(chat2.Severe.UserMapChannel[tid].ChatList[uid], consts.ChatSendTextMsg, socket2.ChatSendTextMsgStruct{
				ID:        msgInfo.ID,
				Uid:       msgInfo.Uid,
				Username:  msgInfo.UInfo.Username,
				Photo:     photo,
				Tid:       msgInfo.Tid,
				Message:   msgInfo.Message,
				Type:      msgInfo.Type,
				CreatedAt: msgInfo.CreatedAt,
			})
			return
		} else {
			// tid用户不处在与uid用户的聊天窗口，添加ChatList未读记录并推送，到tid的ChatListWs
			// 找到tid用户私信列表中关于uid的聊天记录
			cl := new(chat.ChatsListInfo)
			err := cl.UnreadAutoCorrection(tid, uid)
			if err != nil {
				global.Logger.Errorf("uid:%d tid:%d 私信列表未读信息记录自增失败", tid, uid)
			}
			ci := new(chat.ChatsListInfo)
			_ = ci.FindByID(uid, tid)
			// 推送到tid的私信列表
			response.SuccessWs(chat2.Severe.UserMapChannel[tid].Socket, consts.ChatUnreadNotice, socket2.ChatUnreadNoticeStruct{
				Uid:         uid,
				Tid:         tid,
				LastMessage: ci.LastMessage,
				LastMessageInfo: socket2.ChatSendTextMsgStruct{
					ID:        msgInfo.ID,
					Uid:       msgInfo.Uid,
					Username:  msgInfo.UInfo.Username,
					Photo:     photo,
					Tid:       msgInfo.Tid,
					Message:   msgInfo.Message,
					Type:      msgInfo.Type,
					CreatedAt: msgInfo.CreatedAt,
				},
				Unread: cl.Unread,
			})
		}
	} else {
		// tid 用户不在线，直接更新数据库中未读数量，等 tid用户上线后进入私信列表界面 ChatWs 会自动推送私信列表未读信息数量
		cl := new(chat.ChatsListInfo)
		err := cl.UnreadAutoCorrection(tid, uid)
		if err != nil {
			global.Logger.Errorf("uid:%d tid:%d 私信列表未读信息记录自增失败", tid, uid)
		}
	}
}
