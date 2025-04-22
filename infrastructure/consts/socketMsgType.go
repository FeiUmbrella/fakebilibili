package consts

const (
	// VideoSocketTypeError ws视频推送错误
	VideoSocketTypeError = "error"
	// VideoSocketTypeNumber 返回在线观看人数
	VideoSocketTypeNumberOfViewers = "numberOfViewers"
	// VideoSocketTypeSendBarrage 发送视频弹幕，并通知视频房间里的所有人
	VideoSocketTypeSendBarrage        = "sendBarrage"
	VideoSocketTypeResponseBarrageNum = "responseBarrageNum"
	// NoticeSocketTypeMessage 消息通知
	NoticeSocketTypeMessage = "messageNotice"
	// ChatSendTextMsg 聊天界面推送文本消息
	ChatSendTextMsg = "chatSendTextMsg"
	// ChatUnreadNotice 聊天消息未读通知
	ChatUnreadNotice = "chatUnreadNotice"
	// ChatOnlineUnreadMsg 用户上线后推送未读私信
	ChatOnlineUnreadMsg = "chatOnlineUnreadMsg"

	// ************ live Socket 相关 *****************//
	// Error 错误信息
	Error = "error"
	/*
		WebClientBarrageReq  发送弹幕请求数据
		WebClientBarrageRes  发送弹幕响应数据
		WebClientHistoricalBarrageRes 历史弹幕消息
	*/
	WebClientBarrageReq           = "webClientBarrageReq"
	WebClientBarrageRes           = "webClientBarrageRes"
	WebClientHistoricalBarrageRes = "webClientHistoricalBarrageRes"

	//WebClientEnterLiveRoomRes  用户上下线提醒
	WebClientEnterLiveRoomRes = "webClientEnterLiveRoomRes"
)
