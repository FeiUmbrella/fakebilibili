package captcha

import (
	"fakebilibili/infrastructure/pkg/global"
	"github.com/dchest/captcha"
)

type captchaResponse struct {
	CaptchaId string `json:"captchaId"`
	ImageUrl  string `json:"imageUrl"`
}

// todo:这里没有完全实现验证图片功能，只有个GetCaptcha
func GetCaptcha(id string) (results interface{}, err error) {
	length := captcha.DefaultLen
	captchaId := captcha.NewLen(length)
	var response = &captchaResponse{
		CaptchaId: captchaId,
		ImageUrl:  "/captcha/" + captchaId + ".png",
	}
	global.Logger.Infof("captcha信息为%s,%s", response.CaptchaId, response.ImageUrl)
	return response, nil
}
