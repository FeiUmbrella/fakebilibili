package validator

import (
	"fmt"
	"regexp"
)

// VerifyMobileFormat 手机号格式验证
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// VerifyEmailFormat 邮箱格式验证
func VerifyEmailFormat(email string) bool {
	regex := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	return regex.MatchString(email)
}

// CheckVideoSuffix 检查要保存的文件后缀是否合法
func CheckVideoSuffix(suffix string) error {
	switch suffix {
	case ".jpg", ".jpeg", ".png", ".ico", ".gif", ".wbmp", ".bmp", ".svg", ".webp", ".mp4":
		return nil
	default:
		return fmt.Errorf("非法后缀！")
	}
}
