package live

import (
	"fakebilibili/domain/servicce/live"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/pkg/utils/response"
	"fakebilibili/infrastructure/proto/pb"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"strconv"
)

func (lv LivesControllers) LiveSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)

	// 判断是否为该用户创建了直播间
	liveRoom, _ := strconv.Atoi(ctx.Query("liveRoom"))
	liveRoomID := uint(liveRoom)
	if live.Severe.LiveRoom[liveRoomID] == nil {
		// 该直播间不存在
		message := &pb.Message{
			MsgType: consts.Error,
			Data:    []byte("直播未开始"),
		}
		res, _ := proto.Marshal(message)
		_ = ws.WriteMessage(websocket.BinaryMessage, res) // ws传输二进制信息
		return
	}
	// 直播间存在
	err := live.CreateSocket(ctx, uid, liveRoomID, ws)
	if err != nil {
		response.ErrorWs(ws, err.Error())
		return

	}
}
