package router

import (
	"fakebilibili/adapter/http/router/users"
	"github.com/gin-gonic/gin"
)

type RoutersGroup struct {
	Users users.RouterGroup
}

var RoutersGroupApp = new(RoutersGroup)

func InitRouter() {
	router := gin.Default()
	//todo:跨域中间件
	//router.Use(middlewares.Cors())
	PrivateGroup := router.Group("")
	PrivateGroup.Use()
	{
		// 静态资源访问
		router.Static("/assets", "./assets")
		RoutersGroupApp.Users.LoginRouter.InitLoginRouter(PrivateGroup)
	}
	err := router.Run(":8081")
	if err != nil {
		return
	}
}
