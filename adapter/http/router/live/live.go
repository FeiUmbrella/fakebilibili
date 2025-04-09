package live

import (
	"fakebilibili/adapter/http/controller/live"
	"fakebilibili/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

type LivesRouter struct {
}

func (s *LivesRouter) InitLiveRouter(Router *gin.RouterGroup) {
	liveRouter := Router.Group("live").Use(middleware.VerificationToken())
	{
		liveControllers := new(live.LivesControllers)
		liveRouter.POST("/getLiveRoom", liveControllers.GetLiveRoom)         // 主播进行开播，获取推流地址和推流码
		liveRouter.POST("/getLiveRoomInfo", liveControllers.GetLiveRoomInfo) // 观众进入直播间，返回直播间信息以及该直播间的拉流地址
		liveRouter.POST("/getBeLiveList", liveControllers.GetBeLiveList)     // 返回所有正在直播的直播间列表
	}
}
