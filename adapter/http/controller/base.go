package controller

import (
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type BaseControllers struct{}

// Response 传入逻辑层结果，进一步判断产生出控制器的输出结果
func (c BaseControllers) Response(ctx *gin.Context, results interface{}, err error) {
	if err != nil {
		// 逻辑层执行报错处理
		response.Error(ctx, err.Error())
		return
	}
	// 逻辑层返回正常结果
	response.Success(ctx, results)
}

// ShouldBind 将请求中的参数绑定到对应结构体
func ShouldBind[T interface{}](ctx *gin.Context, data T) (t T, err error) {
	if err = ctx.ShouldBind(data); err != nil {
		global.Logger.Errorf("传入请求中的参数绑定失败，type:%T; 错误原因: %s ", data, err.Error())
		// todo:进一步查看绑定过程哪里出错，快速确定出错位置
		//validator.CheckParams(ctx, err)
		return t, err
	}
	return data, nil
}
