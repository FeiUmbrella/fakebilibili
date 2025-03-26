package middleware

import (
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/jwt"
	response2 "fakebilibili/infrastructure/pkg/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

// VerificationToken 提取请求头中的token并解析，来判断用户是否是登录状态
func VerificationToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if len(token) == 0 {
			response2.NotLogin(c, "未登录，Authorization Token为空")
			c.Abort()
			return
		}
		claim, err := jwt.ParseToken(token)
		if err != nil {
			response2.NotLogin(c, "Token过期")
			c.Abort()
			return
		}
		redisToken, _ := global.RedisDb.Get(fmt.Sprintf("%d_%s", claim.UserID, consts.TokenString)).Result()
		if redisToken != token {
			// 传进来的token与redis保存的不一致，有两种情况：
			//1.传错了，但是实际应用中用户是不用手动输入的，排除这种可能；2.已在别处登录获得新的token，所以旧的token不能用了
			response2.NotLogin(c, "已在别处登录")
			c.Abort()
			return
		}

		u := new(user.User)
		if !u.IsExistByField("id", claim.UserID) {
			// 数据库没有改动的情况下
			response2.NotLogin(c, "用户异常")
			c.Abort()
			return
		}
		c.Set("uid", claim.UserID)
		c.Set("currentUserName", u.Username)
		c.Next()
	}
}

// VerificationTokenNotNecessary 不强制要求用户为登录状态
func VerificationTokenNotNecessary() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if len(token) == 0 {
			// 不强制要求用户为登录状态
			c.Next()
		} else { // 用户登录状态
			claim, err := jwt.ParseToken(token)
			if err != nil {
				c.Next()
			}
			u := new(user.User)
			if !u.IsExistByField("id", claim.UserID) {
				// 数据库没有改动的情况下
				response2.NotLogin(c, "用户异常")
				c.Abort()
				return
			}
			c.Set("uid", claim.UserID)
			c.Set("currentUserName", u.Username)
			c.Next()
		}
	}
}
