package users

import (
	"fakebilibili/domain/servicce/users"
	"fakebilibili/domain/servicce/users/chat"
	"fakebilibili/domain/servicce/users/chatUser"
	"fakebilibili/infrastructure/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
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

// ChatByUserSocket 对前端用户在线时跟其他用户聊天接收到的聊天信息 进行推送
func (us UserControllers) ChatByUserSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	tidQuery, _ := strconv.Atoi(ctx.Query("tid"))
	tid := uint(tidQuery)
	ws := conn.(*websocket.Conn)
	err := chatUser.CreateChatByUserSocket(uid, tid, ws)
	if err != nil {
		response.ErrorWs(ws, "创建用户聊天socket失败")
	}
}
