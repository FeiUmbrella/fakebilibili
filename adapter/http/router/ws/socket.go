package ws

import (
	"fakebilibili/adapter/http/controller/contribution"
	"fakebilibili/adapter/http/controller/live"
	"fakebilibili/adapter/http/controller/users"
	"fakebilibili/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *RouterGroup) InitSocketRouter(Router *gin.RouterGroup) {
	socketRouter := Router.Group("ws").Use(middleware.VerificationTokenAsSocket())
	{
		userControllers := new(users.UserControllers)
		liveControllers := new(live.LivesControllers)
		contributionControllers := new(contribution.Controllers)
		socketRouter.GET("/noticeSocket", userControllers.NoticeSocket)
		socketRouter.GET("/chatSocket", userControllers.ChatSocket)
		socketRouter.GET("/chatUserSocket", userControllers.ChatByUserSocket)
		socketRouter.GET("/liveSocket", liveControllers.LiveSocket)
		socketRouter.GET("/videoSocket", contributionControllers.VideoSocket)
	}
}
