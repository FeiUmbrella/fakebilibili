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

	// UniqueVideoRecommendPrefix 将按照热度选出放在主页的视频id保存在bitmap
	UniqueVideoRecommendPrefix = "uniqueVideoRecommendPrefix_"

	//查询视频相关信息
	VideoBarragePrefix     = "videoBarrageOf_"         //查询视频的弹幕信息
	VideoCommentZSetPrefix = "VideoCommentZSetPrefix_" //zset中一个视频的key
	VideoCommentHashPrefix = "VideoCommentHashPrefix_" // hash中一个视频的key

	//VideoWatchByID 观看视频
	VideoWatchByID   = "videoWatchBy_"
	ArticleWatchByID = "articleWatchBy_" // 观看文章

	// 推荐视频列表信息
	RecommendVideosList = "RecommendVideosList"
)
