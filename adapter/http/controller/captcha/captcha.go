package captcha

import (
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/domain/servicce/captcha"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	controller.BaseControllers
}

func (c *Controllers) GetCaptcha(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.GetCaptchaStruct)); err == nil {
		res, err := captcha.GetCaptcha(rec.CaptchaId)
		c.Response(ctx, res, err)
	}
}
