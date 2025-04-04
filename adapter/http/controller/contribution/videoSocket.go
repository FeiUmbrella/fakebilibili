package contribution

import (
	"fakebilibili/domain/servicce/contribution/socket"
	"fakebilibili/infrastructure/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
)

// VideoSocket  向观看某一视频的用户推送观看该视频的人数
func (c Controllers) VideoSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	//判断是否创建视频socket房间
	id, _ := strconv.Atoi(ctx.Query("videoID"))
	videoID := uint(id)
	//无人观看主动创建
	if socket.Severe.VideoRoom[videoID] == nil {
		socket.Severe.VideoRoom[videoID] = make(socket.UserMapChannel, 10)
	}
	err := socket.CreateVideoSocket(uid, videoID, ws)
	if err != nil {
		response.ErrorWs(ws, "创建socket失败")
	}
}
