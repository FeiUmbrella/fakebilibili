package user

import (
	"crypto/md5"
	"fakebilibili/infrastructure/model/user/liveInfo"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
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
	Photo       datatypes.JSON `json:"photo" gorm:"column:photo"` // 头像
	Gender      int8           `json:"gender" gorm:"column:gender"`
	BirthDate   time.Time      `json:"birth_date" gorm:"column:birth_date"`                       // 注册日期
	IsVisible   int8           `json:"is_visible" gorm:"column:is_visible"`                       // todo: 这个字段干什么用的？
	Signature   string         `json:"signature" gorm:"column:signature;type:varchar(255)"`       // todo: 这个字段干什么用的？
	SocialMedia string         `json:"social_media" gorm:"column:social_media;type:varchar(255)"` // todo: 这个字段干什么用的？

	LiveInfo liveInfo.LiveInfo `json:"liveInfo" gorm:"foreignKey:Uid"`
}

type UserList []User

func (User) TableName() string {
	return "lv_users"
}

// dao层 ---------------------------------------

// IsExistByField 查找user表中字段field是否存在value的字段
func (us *User) IsExistByField(field string, value any) bool {
	err := global.MysqlDb.Where(field, value).First(&us).Error
	if err != nil {
		return false
	}
	return true
}

// Create 创建用户
func (us *User) Create() error {
	return global.MysqlDb.Model(&User{}).Create(&us).Error
}

// IfPasswordCorrect 判断用户密码是否正确
func (us *User) IfPasswordCorrect(userPassword string) bool {
	pwd := us.Salt + userPassword + us.Salt
	pwdMd5 := md5.Sum([]byte(pwd))
	pwdMd5Str := fmt.Sprintf("%x", pwdMd5)
	return pwdMd5Str == us.Password
}

// Update 更新用户信息
func (us *User) Update() bool {
	err := global.MysqlDb.Model(&User{}).Where("id = ?", us.ID).Updates(us).Error
	return err == nil
}

// Find 查找用户信息
func (us *User) Find(uid uint) {
	global.MysqlDb.Model(&User{}).Where("id = ?", uid).Find(&us)
}
