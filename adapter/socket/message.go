package socket

import (
	"fakebilibili/domain/servicce/users"
	"fakebilibili/domain/servicce/users/chat"
)

func init() {
	//初始化所有socket
	// 用户上线后向前端用户推送未读通知
	go users.Severe.Start()
	// 前端用户上线后向前端用户推送未读私信
	go chat.Severe.Start()
}
