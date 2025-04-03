package users

import (
	"fakebilibili/domain/servicce/users"
	"fakebilibili/domain/servicce/users/chat"
	"fakebilibili/infrastructure/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// NoticeSocket 推送通知的ws
func (us UserControllers) NoticeSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	err := users.CreateNoticeSocket(uid, ws)
	if err != nil {
		response.ErrorWs(ws, "创建通知socket失败")
	}
}

// ChatSocket 点击私信聊天列表加载离线时接收到的私信
func (us UserControllers) ChatSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	err := chat.CreateChatSocket(uid, ws)
	if err != nil {
		response.ErrorWs(ws, "创建推送未读私信socket失败")
	}
}
