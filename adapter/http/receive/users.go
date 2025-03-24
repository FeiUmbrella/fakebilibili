package receive

// RegisterReceiveStruct 用户注册相关参数结构体
type RegisterReceiveStruct struct {
	UserName         string `json:"username" binding:"required" form:"username"`
	Password         string `json:"password" binding:"required" form:"password"`
	Email            string `json:"email" binding:"required,email" form:"email"`
	VerificationCode string `json:"verificationCode" binding:"required" form:"verificationCode"`
}

// SendEmailVerCodeReceiveStruct 发送邮箱验证码
type SendEmailVerCodeReceiveStruct struct {
	Email string `json:"email" binding:"required,email" form:"email"`
}
