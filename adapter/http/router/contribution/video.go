package contribution

import (
	"fakebilibili/adapter/http/controller/contribution"
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
	}
}
