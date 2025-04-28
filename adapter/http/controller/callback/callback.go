package callback

import (
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive/callback"
	callback2 "fakebilibili/domain/servicce/callback"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	controller.BaseControllers
}

// AliyunTranscodingMedia 阿里云媒体转码成功回调
func (c *Controllers) AliyunTranscodingMedia(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(callback.AliyunMediaCallback[callback.AliyunTranscodingMediaStruct])); err == nil {
		res, err := callback2.AliyunTranscodingMedia(rec)
		c.Response(ctx, res, err)
	}
}
