package router

import (
	"fakebilibili/adapter/http/middleware"
	"fakebilibili/adapter/http/router/OSSCommonality"
	"fakebilibili/adapter/http/router/callback"
	"fakebilibili/adapter/http/router/contribution"
	"fakebilibili/adapter/http/router/home"
	"fakebilibili/adapter/http/router/live"
	"fakebilibili/adapter/http/router/users"
	"fakebilibili/adapter/http/router/ws"
	"github.com/gin-gonic/gin"
)

type RoutersGroup struct {
	Users          users.RouterGroup
	Ws             ws.RouterGroup
	Contribution   contribution.RouterGroup
	Live           live.RouterGroup
	Home           home.RouterGroup
	OSSCommonality osscommonality.RouterGroup
	Callback       callback.RouterGroup
}

var RoutersGroupApp = new(RoutersGroup)

func InitRouter() {
	router := gin.Default()
	router.Use(middleware.Cors())
	PrivateGroup := router.Group("")
	PrivateGroup.Use()
	{
		// 静态资源访问
		router.Static("/assets", "./assets")
		RoutersGroupApp.Users.LoginRouter.InitLoginRouter(PrivateGroup)
		RoutersGroupApp.Users.SpaceRouter.InitSpaceRouter(PrivateGroup)
		RoutersGroupApp.Users.InitRouter(PrivateGroup)
		RoutersGroupApp.Ws.InitSocketRouter(PrivateGroup)
		RoutersGroupApp.Live.InitLiveRouter(PrivateGroup)
		RoutersGroupApp.Home.InitHomeRouter(PrivateGroup)
		RoutersGroupApp.OSSCommonality.InitRouter(PrivateGroup)
		RoutersGroupApp.Contribution.VideoRouter.InitVideoRouter(PrivateGroup)
		RoutersGroupApp.Contribution.ArticleRouter.InitArticleRouter(PrivateGroup)
		RoutersGroupApp.Callback.InitRouter(PrivateGroup)
	}
	err := router.Run(":8081")
	if err != nil {
		return
	}
}
