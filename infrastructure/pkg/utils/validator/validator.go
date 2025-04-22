package validator

import (
	"errors"
	"fakebilibili/infrastructure/pkg/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var (
	ValidTrans ut.Translator
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

// CheckParams 检测解析http请求中数据时发生的错误
func CheckParams(ctx *gin.Context, err error) {
	if err != nil {
		// 传进来的是一个错误
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) { // 该错误是一个参数校验错误
			for _, fieldError := range err.(validator.ValidationErrors) {
				/*
					使用翻译器 ValidTrans 将错误信息翻译成指定语言
					fieldError.Tag(): 获取验证规则标签（如 "required", "email" 等）
					fieldError.Field(): 获取字段名
					fieldError.Param(): 获取验证规则的参数（如最小长度值等）
				*/
				msg, _ := ValidTrans.T(fieldError.Tag(), fieldError.Field(), fieldError.Param())
				response.Error(ctx, msg)
				return
			}
		} else {
			// 是一个其他错误，直接返回原始错误信息
			response.TypeError(ctx, err.Error())
			return
		}
	}
}
