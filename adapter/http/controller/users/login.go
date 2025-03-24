package users

import (
	"errors"
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/domain/servicce/users"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/limiter"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"time"
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
		// 令牌桶来限制邮箱频繁发送 每10s放一个令牌，桶大小为10
		l := limiter.NewLimiter(rate.Every(10*time.Second), 1, rec.Email)
		if !l.Allow() {
			lg.Response(ctx, nil, errors.New("请求过于频繁，请1分钟后再试"))
			return
		}
		results, err := users.SendEmailVerCode(rec)
		global.Logger.Infof("向邮箱：%s发送验证码", rec.Email)
		lg.Response(ctx, results, err)
	}
}
