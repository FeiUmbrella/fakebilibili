package contribution

import (
	"fakebilibili/adapter/http/controller/contribution"
	"fakebilibili/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

type VideoRouter struct {
}

func (v *VideoRouter) InitVideoRouter(Router *gin.RouterGroup) {
	contributionControllers := new(contribution.Controllers)

	// 不需要登入
	contributionRouterNoVerification := Router.Group("contribution").Use()
	{
		contributionRouterNoVerification.GET("/video/barrage/v3/", contributionControllers.GetVideoBarrage)
		contributionRouterNoVerification.GET("/getVideoBarrageList", contributionControllers.GetVideoBarrageList)
		contributionRouterNoVerification.POST("/getVideoComment", contributionControllers.GetVideoComment)
		contributionRouterNoVerification.POST("/getVideoCommentCountById", contributionControllers.GetVideoCommentCountById)
		contributionRouterNoVerification.POST("/likeVideoComment", contributionControllers.LikeVideoComment) //给评论点赞，先不需要登录
	}

	// 非必须登入
	contributionRouterNotNecessary := Router.Group("contribution").Use(middleware.VerificationTokenNotNecessary())
	{
		contributionRouterNotNecessary.POST("/getVideoContributionByID", contributionControllers.GetVideoContributionByID)
	}
	// body中携带token的http请求
	contributionRouterParameter := Router.Group("contribution").Use(middleware.VerificationTokenAsParameter())
	{
		contributionRouterParameter.POST("/video/barrage/v3/", contributionControllers.SendVideoBarrage)
	}

	// 请求头Header中携带token的请求
	contributionRouter := Router.Group("contribution").Use(middleware.VerificationToken())
	{
		contributionRouter.POST("/createVideoContribution", contributionControllers.CreateVideoContribution)
		contributionRouter.POST("/updateVideoContribution", contributionControllers.UpdateVideoContribution)
		contributionRouter.POST("/deleteVideoByID", contributionControllers.DeleteVideoByID)
		contributionRouter.POST("/videoPostComment", contributionControllers.VideoPostComment)
		contributionRouter.POST("/getVideoManagementList", contributionControllers.GetVideoManagementList)
		contributionRouter.POST("/likeVideo", contributionControllers.LikeVideo)
		contributionRouter.POST("/deleteVideoByPath", contributionControllers.DeleteVideoByPath)
		contributionRouter.GET("/getLastWatchTime", contributionControllers.GetLastWatchTime)
		contributionRouter.POST("/sendWatchTime", contributionControllers.SendWatchTime)
	}
}
