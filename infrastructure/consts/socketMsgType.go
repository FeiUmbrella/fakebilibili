package consts

const (
	// VideoSocketTypeError ws视频推送错误
	VideoSocketTypeError = "error"
	// VideoSocketTypeNumber 返回在线观看人数
	VideoSocketTypeNumber = "numberOfViewers"
	// VideoSocketTypeSendBarrage 发送弹幕
	VideoSocketTypeSendBarrage = "sendBarrage"
	// VideoSocketTypeSendBarrage 发送弹幕 todo:这个是啥？
	VideoSocketTypeResponseBarrageNum = "responseBarrageNum"
	// NoticeSocketTypeMessage 消息通知
	NoticeSocketTypeMessage = "messageNotice"
	// ChatSendTextMsg 聊天界面推送文本消息
	ChatSendTextMsg = "chatSendTextMsg"
	// ChatUnreadNotice 聊天消息未读通知
	ChatUnreadNotice = "chatUnreadNotice"
	// todo:这个表示什么意思，怎么chat又notice？
	ChatOnlineUnreadNotice = "chatOnlineUnreadNotice"
)
