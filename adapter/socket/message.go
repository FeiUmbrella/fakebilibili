package socket

import (
	"fakebilibili/domain/servicce/live"
	"fakebilibili/domain/servicce/users"
	"fakebilibili/domain/servicce/users/chat"
	"fakebilibili/domain/servicce/users/chatUser"
)

func init() {
	//初始化所有socket
	// 用户上线后向前端用户推送未读通知
	go users.Severe.Start()
	// 前端用户上线后向前端用户推送未读私信
	go chat.Severe.Start()
	// 前端用户向另一用户发送信息，对另一用户进行推送
	go chatUser.Severe.Start()
	// 前端用户进入直播间，发送弹幕推送给直播间其他用户、向新进入直播间的用户推送历史弹幕、用户退出/进入直播间时向直播间其他用户广播
	go live.Severe.Start()
}
