package receive

// RegisterReceiveStruct 用户注册相关参数结构体
type RegisterReceiveStruct struct {
	UserName         string `json:"username" binding:"required" form:"username"`
	Password         string `json:"password" binding:"required" form:"password"`
	Email            string `json:"email" binding:"required,email" form:"email"`
	VerificationCode string `json:"verificationCode" binding:"required" form:"verificationCode"` // 邮箱验证码
}

// SendEmailVerCodeReceiveStruct 发送邮箱验证码
type SendEmailVerCodeReceiveStruct struct {
	Email string `json:"email" binding:"required,email" form:"email"`
}

// LoginReceiveStruct 用户登录 相关参数结构体
type LoginReceiveStruct struct {
	UserName string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
}

// ForgetReceiveStruct 忘记密码进行找回 相关参数结构体
type ForgetReceiveStruct struct {
	Password         string `json:"password" binding:"required" form:"password"`
	Email            string `json:"email" binding:"required,email" form:"email"`
	VerificationCode string `json:"verificationCode" binding:"required" form:"verificationCode"` // 邮箱验证码
}

// GetAttentionListReceiveStruct 获取用户关注列表 相关参数结构体
type GetAttentionListReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// GetVermicelliListReceiveStruct 获取用户粉丝列表 相关参数结构体
type GetVermicelliListReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// GetSpaceIndividualReceiveStruct 获取用户个人空间 相关参数结构体
type GetSpaceIndividualReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// GetReleaseInformationReceiveStruct 获取用户发布的作品(视频and专栏) 相关参数结构体
type GetReleaseInformationReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}
