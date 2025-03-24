package response

import (
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"time"
)

// UserInfoResponseStruct 返回前端用户信息结构体
type UserInfoResponseStruct struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"username"`
	Photo     string    `json:"photo"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

// UserInfoResponse  生成返回用用户信息结构体
func UserInfoResponse(us *user.User, token string) UserInfoResponseStruct {
	//判断用户是否为微信用户进行图片处理
	photo, _ := conversion.FormattingJsonSrc(us.Photo)
	return UserInfoResponseStruct{
		ID:        us.ID,
		UserName:  us.Username,
		Photo:     photo,
		Token:     token,
		CreatedAt: us.CreatedAt,
	}
}
