package home

import (
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive/home"
	home2 "fakebilibili/domain/servicce/home"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	controller.BaseControllers
}

// GetHomeInfo 获取主页信息-轮播图和推荐视频
func (c Controllers) GetHomeInfo(ctx *gin.Context) {
	// 参数有page、size、keyword
	if rec, err := controller.ShouldBind(ctx, new(home.GetHomeInfoReceiveStruct)); err == nil {
		results, err := home2.GetHomeInfo(rec)
		c.Response(ctx, results, err)
	}
}

// SubmitBug 用户提交bug
func (c Controllers) SubmitBug(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(home.SubmitBugReceiveStruct)); err == nil {
		results, err := home2.SubmitBug(rec)
		c.Response(ctx, results, err)
	}
}
