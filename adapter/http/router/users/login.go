package users

import (
	"fakebilibili/adapter/http/controller/users"
	"github.com/gin-gonic/gin"
)

type LoginRouter struct {
}

func (s *LoginRouter) InitLoginRouter(Router *gin.RouterGroup) {
	loginRouter := Router.Group("login").Use()
	{
		loginControllers := new(users.LoginControllers)
		loginRouter.POST("/sendEmailVerificationCode", loginControllers.SendEmailVerCode) // 注册发送邮箱验证码
		loginRouter.POST("/register", loginControllers.Register)                          // 注册
		loginRouter.POST("/login", loginControllers.Login)
		loginRouter.POST("/sendEmailVerificationCodeByForget", loginControllers.SendEmailVerCodeByForget)
		loginRouter.POST("/forget", loginControllers.Forget)
	}
}
