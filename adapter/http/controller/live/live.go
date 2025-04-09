package live

import (
	"fakebilibili/adapter/http/controller"
	live2 "fakebilibili/adapter/http/receive/live"
	"fakebilibili/domain/servicce/live"
	"github.com/gin-gonic/gin"
)

type LivesControllers struct {
	controller.BaseControllers
}

// GetLiveRoom 返回开播时对应直播间推流地址和推流码
func (lv LivesControllers) GetLiveRoom(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	res, err := live.GetLiveRoom(uid)
	lv.Response(ctx, res, err)
}

// GetLiveRoomInfo 给前端返回直播间信息及直播间拉流地址
func (lv LivesControllers) GetLiveRoomInfo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(live2.GetLiveRoomInfoReceiveStruct)); err == nil {
		result, err := live.GetLiveRoomInfo(*rec, uid)
		lv.Response(ctx, result, err)
	}
}

// GetBeLiveList 获取开通直播的所有用户
func (lv LivesControllers) GetBeLiveList(ctx *gin.Context) {
	res, err := live.GetBeLiveList()
	lv.Response(ctx, res, err)
}
