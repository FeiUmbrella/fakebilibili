package home

import (
	"fakebilibili/adapter/http/controller/home"
	"github.com/gin-gonic/gin"
)

type homeRouter struct {
}

func (s *homeRouter) InitHomeRouter(Router *gin.RouterGroup) {
	homeRouter := Router.Group("home")
	{
		homeControllers := new(home.Controllers)
		homeRouter.POST("/getHomeInfo", homeControllers.GetHomeInfo)
		homeRouter.POST("/submitBug", homeControllers.SubmitBug)
	}
}
