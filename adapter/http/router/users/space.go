package users

import (
	"fakebilibili/adapter/http/controller/users"
	"fakebilibili/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

type SpaceRouter struct {
}

func (s *SpaceRouter) InitSpaceRouter(Router *gin.RouterGroup) {
	spaceControllers := new(users.SpaceControllers)

	// 必须要登录，用中间件检验是否登录
	SpaceRouter1 := Router.Group("space").Use(middleware.VerificationToken())
	{
		// 获取关注列表
		SpaceRouter1.POST("/getAttentionList", spaceControllers.GetAttentionList)
		// 获取粉丝列表
		SpaceRouter1.POST("/getVermicelliList", spaceControllers.GetVermicelliList)
	}

	// 非必须登录
	SpaceRouter2 := Router.Group("space").Use(middleware.VerificationTokenNotNecessary())
	{
		// 获取个人空间
		SpaceRouter2.POST("/getSpaceIndividual", spaceControllers.GetSpaceIndividual)
		// 获取
		SpaceRouter2.POST("/getReleaseInformation", spaceControllers.GetReleaseInformation)
	}
}
