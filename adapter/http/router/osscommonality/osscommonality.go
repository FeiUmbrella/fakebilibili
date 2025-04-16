package osscommonality

import (
	"fakebilibili/adapter/http/controller/OSSCommonality"
	"fakebilibili/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *RouterGroup) InitRouter(Router *gin.RouterGroup) {
	ossCommonalityControllers := new(osscommonality.Controllers)
	routers := Router.Group("commonality").Use()
	{
		routers.POST("/ossSTS", ossCommonalityControllers.OssSTS)
		routers.POST("/upload", ossCommonalityControllers.Upload)
		routers.POST("/UploadSlice", ossCommonalityControllers.UploadSlice)
		routers.POST("/uploadCheck", ossCommonalityControllers.UploadCheck)
		routers.POST("/uploadMerge", ossCommonalityControllers.UploadMerge)
		routers.POST("/uploadingMethod", ossCommonalityControllers.UploadingMethod)
		routers.POST("/uploadingDir", ossCommonalityControllers.UploadingDir)
		routers.POST("/getFullPathOfImage", ossCommonalityControllers.GetFullPathOfImage)
		routers.POST("/registerMedia", ossCommonalityControllers.RegisterMedia)
	}

	// 非必须登入
	router := Router.Group("commonality").Use(middleware.VerificationTokenNotNecessary())
	{
		router.POST("/search", ossCommonalityControllers.Search)
	}
}
