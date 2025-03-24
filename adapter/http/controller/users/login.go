package users

import (
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/domain/servicce/users"
	"github.com/gin-gonic/gin"
)

type LoginControllers struct {
	controller.BaseControllers
}

// Register 用户注册
func (lg LoginControllers) Register(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.RegisterReceiveStruct)); err == nil {
		results, err := users.Register(rec)
		// 将逻辑层的处理结果出入lg.Response，进一步判断
		lg.Response(ctx, results, err)
	}
}

// SendEmailVerCode 用户注册发送邮箱验证码
func (lg LoginControllers) SendEmailVerCode(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.SendEmailVerCodeReceiveStruct)); err == nil {
		// todo:令牌桶来限制邮箱频繁发送

		results, err := users.SendEmailVerCode(rec)
		lg.Response(ctx, results, err)
	}
}
