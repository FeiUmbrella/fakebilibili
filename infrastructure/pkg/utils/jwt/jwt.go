package jwt

import (
	"errors"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// SaltStr  生成密码盐的随机字符串
var SaltStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// HotKey 密钥
var HotKey = []byte("fakebilibili")

type Claims struct {
	UserID    uint
	LoginTime string
	jwt.StandardClaims
}

// GenerateToken 生成Auth Token
func GenerateToken(uid uint) string {
	expireTime := time.Now().Add(7 * 24 * time.Hour).Unix()

	calims := &Claims{
		UserID:    uid,
		LoginTime: time.Now().Format("2006-01-02 15:04:05"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime, // 过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "root",       // 颁发者
			Subject:   "Auth Token", // 签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, calims)
	tokenString, err := token.SignedString(HotKey)
	if err != nil {
		fmt.Println(err)
	}

	// 生成的 Auth Token 向redis中放一份
	key := fmt.Sprintf("%d_%s", uid, consts.TokenString)
	global.RedisDb.Set(key, tokenString, 7*24*time.Hour)
	global.Logger.Infof("给用户 %d 签发并向redis投递了token：%s", uid, tokenString)
	return tokenString
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return HotKey, nil
	})
	if err != nil {
		global.Logger.Errorf("Auth Token parse err: %s", err.Error())
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
