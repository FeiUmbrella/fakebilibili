package user

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// User 基本用户信息
type User struct {
	gorm.Model
	Email       string         `json:"email" gorm:"column:email; type:varchar(255)"`
	Username    string         `json:"username" gorm:"column:username; type:varchar(255)"`
	Openid      string         `json:"openid" gorm:"column:openid; type:varchar(255)"` // 用于微信登录
	Salt        string         `json:"salt" gorm:"column:salt;type:varchar(255)"`      // 加密盐
	Password    string         `json:"password" gorm:"column:password; type:varchar(255)"`
	Photo       datatypes.JSON `json:"photo" gorm:"column:photo"`
	Gender      int8           `json:"gender" gorm:"column:gender"`
	BirthDate   time.Time      `json:"birth_date" gorm:"column:birth_date"`
	IsVisible   int8           `json:"is_visible" gorm:"column:is_visible"`                       // todo: 这个字段干什么用的？
	Signature   string         `json:"signature" gorm:"column:signature;type:varchar(255)"`       // todo: 这个字段干什么用的？
	SocialMedia string         `json:"social_media" gorm:"column:social_media;type:varchar(255)"` // todo: 这个字段干什么用的？
}

type UserList []User

func (User) TableName() string {
	return "lv_users"
}
