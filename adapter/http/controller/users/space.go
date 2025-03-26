package users

import (
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/domain/servicce/users"
	"github.com/gin-gonic/gin"
)

type SpaceControllers struct {
	controller.BaseControllers
}

// GetAttentionList 获取关注列表
func (sp SpaceControllers) GetAttentionList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.GetAttentionListReceiveStruct)); err == nil {
		res, err := users.GetAttentionList(rec, uid)
		sp.Response(ctx, res, err)
	}
}

// GetVermicelliList 获取粉丝列表
func (sp SpaceControllers) GetVermicelliList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.GetVermicelliListReceiveStruct)); err == nil {
		res, err := users.GetVermicelliList(rec, uid)
		sp.Response(ctx, res, err)
	}
}

// GetSpaceIndividual 获取个人空间
func (sp SpaceControllers) GetSpaceIndividual(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.GetSpaceIndividualReceiveStruct)); err == nil {
		res, err := users.GetSpaceIndividual(rec, uid)
		sp.Response(ctx, res, err)
	}
}

// GetReleaseInformation 获取发布信息(视频and专栏)
func (sp SpaceControllers) GetReleaseInformation(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.GetReleaseInformationReceiveStruct)); err == nil {
		res, err := users.GetReleaseInformation(rec)
		sp.Response(ctx, res, err)
	}
}
