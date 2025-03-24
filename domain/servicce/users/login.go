package users

import (
	"crypto/md5"
	"encoding/json"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/adapter/http/response"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/common"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/email"
	"fakebilibili/infrastructure/pkg/utils/jwt"
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

// 默认头像路径
var default_images = []string{consts.DEFAULT_IMAGE_1, consts.DEFAULT_IMAGE_2, consts.DEFAULT_IMAGE_3, consts.DEFAULT_IMAGE_4}

// Register 用户注册相关逻辑
func Register(data *receive.RegisterReceiveStruct) (interface{}, error) {
	users := new(user.User)
	// 检查邮箱
	if users.IsExistByField("email", data.Email) {
		return nil, fmt.Errorf("邮箱已注册")
	}

	//判断验证码是否正确/未到期
	// todo:如果redis key过期，会返回什么？
	verCode, err := global.RedisDb.Get(fmt.Sprintf("%s%s", consts.RegEmailVerCode, data.Email)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("验证码过期")
	}
	if verCode != data.VerificationCode {
		return nil, fmt.Errorf("验证码错误")
	}

	// 密码加密
	// 生成6位密码盐
	salt := make([]byte, 6)
	for i := range salt {
		salt[i] = jwt.SaltStr[rand.Int63()%int64(len(jwt.SaltStr))]
	}
	// p = salt + password + salt
	password := []byte(fmt.Sprintf("%s%s%s", salt, data.Password, salt))
	// md5 加密
	passwordMd5 := fmt.Sprintf("%x", md5.Sum(password))

	// 头像从默认头像随机选一个
	photo, _ := json.Marshal(common.Img{
		Src: getRandomImage(default_images),
		Tp:  "oss",
	})

	// 数据库创建用户
	registerUserData := &user.User{
		Email:     data.Email,
		Username:  data.UserName,
		Salt:      string(salt),
		Password:  passwordMd5,
		Photo:     photo,
		BirthDate: time.Now(),
	}
	err = registerUserData.Create()
	if err != nil {
		return nil, fmt.Errorf("注册失败")
	}

	// 生成对应token
	tokenString := jwt.GenerateToken(registerUserData.ID)
	results := response.UserInfoResponse(registerUserData, tokenString)
	// todo: notice 对应逻辑
	//ne := new(noticeModel.Notice)
	//ne.AddNotice(registerData.ID, 37, 0, noticeModel.UserLogin, "欢迎来到本站，请尽情探索，有任何问题都可以联系我。")
	return results, nil
}

// getRandomImage 随机返回四张图片中的一张作为头像
func getRandomImage(images []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(images))
	return images[index]
}

// SendEmailVerCode 注册时发送邮箱验证码
func SendEmailVerCode(data *receive.SendEmailVerCodeReceiveStruct) (interface{}, error) {
	users := new(user.User)
	if users.IsExistByField("email", data.Email) {
		return nil, fmt.Errorf("邮箱已注册")
	}

	// 发送对象
	mailTo := []string{data.Email}
	// 邮箱主题
	subject := "验证码"
	// 生成6位验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000))
	// 邮件正文
	body := fmt.Sprintf("您的注册验证码为:%s,5分钟内有效,请勿转发他人", code)
	err := email.SendEmail(mailTo, subject, body)
	if err != nil {
		return nil, err
	}
	err = global.RedisDb.Set(fmt.Sprintf("%s%s", consts.RegEmailVerCode, data.Email), code, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}
	return "邮箱验证码发送成功", nil
}
