package router

import (
	"fakebilibili/adapter/http/middleware"
	"fakebilibili/adapter/http/router/users"
	"github.com/gin-gonic/gin"
)

type RoutersGroup struct {
	Users users.RouterGroup
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
	}
	err := router.Run(":8081")
	if err != nil {
		return
	}
}
