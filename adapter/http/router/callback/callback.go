package callback

import (
	"fakebilibili/adapter/http/controller/callback"
	"github.com/gin-gonic/gin"
)

func (r *RouterGroup) InitRouter(Router *gin.RouterGroup) {
	callbackControllers := new(callback.Controllers)
	routers := Router.Group("callback").Use()
	{
		routers.POST("/aliyunTranscodingMedia", callbackControllers.AliyunTranscodingMedia())
	}
}
