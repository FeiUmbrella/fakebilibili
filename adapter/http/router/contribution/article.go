package contribution

import (
	"fakebilibili/adapter/http/controller/contribution"
	"fakebilibili/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

type ArticleRouter struct {
}

func (v *ArticleRouter) InitArticleRouter(Router *gin.RouterGroup) {
	contributionControllers := new(contribution.Controllers)
	// 不需要登录
	contributionRouterNoVerification := Router.Group("contribution").Use()
	{
		contributionRouterNoVerification.POST("/getArticleContributionList", contributionControllers.GetArticleContributionList)
		contributionRouterNoVerification.POST("/getArticleContributionListByUser", contributionControllers.GetArticleContributionListByUser)
		contributionRouterNoVerification.POST("/getArticleComment", contributionControllers.GetArticleComment)
		contributionRouterNoVerification.POST("/getArticleClassificationList", contributionControllers.GetArticleClassificationList)
		contributionRouterNoVerification.POST("/getArticleTotalInfo", contributionControllers.GetArticleTotalInfo)
	}
	// 非必须登入
	contributionRouterNotNecessary := Router.Group("contribution").Use(middleware.VerificationTokenNotNecessary())
	{
		contributionRouterNotNecessary.POST("/getArticleContributionByID", contributionControllers.GetArticleContributionByID)
	}
	// 需要登入
	contributionRouter := Router.Group("contribution").Use(middleware.VerificationToken())
	{
		contributionRouter.POST("/createArticleContribution", contributionControllers.CreateArticleContribution)
		contributionRouter.POST("/updateArticleContribution", contributionControllers.UpdateArticleContribution)
		contributionRouter.POST("/deleteArticleByID", contributionControllers.DeleteArticleByID)
		contributionRouter.POST("/articlePostComment", contributionControllers.ArticlePostComment)
		contributionRouter.POST("/getArticleManagementList", contributionControllers.GetArticleManagementList)
		contributionRouter.POST("/getColumnByClassificationId", contributionControllers.GetColumnByClassificationId)
	}
}
