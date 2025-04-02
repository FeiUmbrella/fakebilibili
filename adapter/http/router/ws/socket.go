package ws

import (
	"fakebilibili/adapter/http/controller/users"
	"fakebilibili/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *RouterGroup) InitSocketRouter(Router *gin.RouterGroup) {
	socketRouter := Router.Group("ws").Use(middleware.VerificationTokenAsSocket())
	{
		userControllers := new(users.UserControllers)
		socketRouter.GET("/noticeSocket", userControllers.NoticeSocket)
	}
}
