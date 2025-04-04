package consts

const (
	/*
		RegEmailVerCode	注册验证码
		RegEmailVerCodeByForget 找回密码验证码
		EmailVerificationCodeByChangePassword 修改密码验证码
	*/
	RegEmailVerCode                       = "regEmailVerCode"
	RegEmailVerCodeByForget               = "regEmailVerCodeByForget"
	EmailVerificationCodeByChangePassword = "emailVerificationCodeByChangePassword"

	// TokenString 用户的Auth Token后缀
	TokenString = "tokenString"

	// LiveRoomHistoricalBarrage 近期的历史弹幕存入redis
	LiveRoomHistoricalBarrage = "liveRoomHistoricalBarrage_"
)
