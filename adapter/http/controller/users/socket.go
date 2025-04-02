package users

import (
	"fakebilibili/domain/servicce/users"
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
