package user

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email       string         `json:"email" gorm:"column:email"`
	Username    string         `json:"username" gorm:"column:username"`
	Openid      string         `json:"openid" gorm:"column:openid"`
	Salt        string         `json:"salt" gorm:"column:salt"`
	Password    string         `json:"password" gorm:"column:password"`
	Photo       datatypes.JSON `json:"photo" gorm:"column:photo"`
	Gender      int8           `json:"gender" gorm:"column:gender"`
	BirthDate   time.Time      `json:"birth_date" gorm:"column:birth_date"`
	IsVisible   int8           `json:"is_visible" gorm:"column:is_visible"`
	Signature   string         `json:"signature" gorm:"column:signature"`
	SocialMedia string         `json:"social_media" gorm:"column:social_media"`
}

type UserList []User

func (User) TableName() string {
	return "lv_users"
}
