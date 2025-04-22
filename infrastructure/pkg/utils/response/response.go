package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Data Controller返回的结果
type Data struct {
	Code    MyCode      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty 当字段为空时，不显示该字段
	Version interface{} `json:"version,omitempty"`
}

func Error(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, &Data{
		Code:    CodeServerBusy,
		Message: msg,
	})
}

func TypeError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, &Data{
		Code:    CodeTypeError,
		Message: CodeTypeError.Msg(),
	})
}

func Default(ctx *gin.Context) {
	rd := &Data{
		Code:    CodeDefault,
		Message: CodeDefault.Msg(),
	}
	ctx.JSON(http.StatusOK, rd)
}

func Success(ctx *gin.Context, data interface{}) {
	rd := &Data{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
}

// NotLogin 返回未登录消息体
func NotLogin(ctx *gin.Context, msg string) {
	rd := &Data{
		Code:    CodeNotLogin,
		Message: msg,
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
}

// BarrageSuccess 弹幕播放器响应
func BarrageSuccess(ctx *gin.Context, data interface{}) {
	rd := &Data{
		Code:    0,
		Message: CodeSuccess.Msg(),
		Data:    data,
		Version: 3,
	}
	ctx.JSON(http.StatusOK, rd)
}
