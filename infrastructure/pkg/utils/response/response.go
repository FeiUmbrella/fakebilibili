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
