package users

import (
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive"
	"fakebilibili/domain/servicce/users"
	"fakebilibili/domain/servicce/users/checkin"
	"github.com/gin-gonic/gin"
)

type UserControllers struct {
	controller.BaseControllers
}

// GetUserInfo 获取用户信息
func (us UserControllers) GetUserInfo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	rec, err := users.GetUserInfo(uid)
	us.Response(ctx, rec, err)
}

// SetUserInfo 设置用户信息
func (us UserControllers) SetUserInfo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.SetUserInfoReceiveStruct)); err == nil {
		res, err := users.SetUserInfo(rec, uid)
		us.Response(ctx, res, err)
	}
}

// DetermineNameExists 判断名字是否存在（应该是在改名或注册时判断）
func (us UserControllers) DetermineNameExists(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.DetermineNameExistsStruct)); err == nil {
		res, err := users.DetermineNameExists(rec, uid)
		us.Response(ctx, res, err)
	}
}

// UpdateAvatar 更新头像
func (us UserControllers) UpdateAvatar(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.UpdateAvatarStruct)); err == nil {
		res, err := users.UpdateAvatar(rec, uid)
		us.Response(ctx, res, err)
	}
}

// GetLiveInfo 获取直播间信息
func (us UserControllers) GetLiveInfo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	res, err := users.GetLiveInfo(uid)
	us.Response(ctx, res, err)
}

// SaveLiveInfo 修改直播间信息
func (us UserControllers) SaveLiveInfo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.SaveLiveDataReceiveStruct)); err == nil {
		res, err := users.SaveLiveInfo(rec, uid)
		us.Response(ctx, res, err)
	}
}

// SendEmailVerificationCodeByChangePassword 在登陆状态下修改密码时发送邮箱验证码
func (us UserControllers) SendEmailVerificationCodeByChangePassword(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	res, err := users.SendEmailVerificationCodeByChangePassword(uid)
	us.Response(ctx, res, err)
}

// ChangePassword 在登陆状态下修改密码
func (us UserControllers) ChangePassword(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.ChangePasswordReceiveStruct)); err == nil {
		results, err := users.ChangePassword(rec, uid)
		us.Response(ctx, results, err)
	}
}

// Attention 关注用户/取消关注
func (us UserControllers) Attention(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.AttentionReceiveStruct)); err == nil {
		results, err := users.Attention(rec, uid)
		us.Response(ctx, results, err)
	}
}

// CreateFavorites 创建或更新收藏夹
func (us UserControllers) CreateFavorites(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.CreateFavoritesReceiveStruct)); err == nil {
		results, err := users.CreateFavorites(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetFavoritesList 以列表形式获取所有收藏夹
func (us UserControllers) GetFavoritesList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := users.GetFavoritesList(uid)
	us.Response(ctx, results, err)
}

// DeleteFavorites 删除收藏夹
func (us UserControllers) DeleteFavorites(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.DeleteFavoritesReceiveStruct)); err == nil {
		results, err := users.DeleteFavorites(rec, uid)
		us.Response(ctx, results, err)
	}
}

// FavoriteVideo 收藏/取消收藏视频
func (us UserControllers) FavoriteVideo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.FavoriteVideoReceiveStruct)); err == nil {
		results, err := users.FavoriteVideo(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetFavoritesListByFavoriteVideo 获取用户中包含该视频的收藏夹列表
func (us UserControllers) GetFavoritesListByFavoriteVideo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.GetFavoritesListByFavoriteVideoReceiveStruct)); err == nil {
		results, err := users.GetFavoritesListByFavoriteVideo(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetFavoriteVideoList 获取收藏夹中的视频列表
func (us UserControllers) GetFavoriteVideoList(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.GetFavoriteVideoListReceiveStruct)); err == nil {
		results, err := users.GetFavoriteVideoList(rec)
		us.Response(ctx, results, err)
	}
}

// GetCollectListName 根据收藏夹id获取收藏夹的title
func (us UserControllers) GetCollectListName(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.GetCollectListNameReceiveStruct)); err == nil {
		results, err := users.GetCollectListName(rec)
		us.Response(ctx, results, err)
	}
}

// GetRecordList 获取历史记录
func (us UserControllers) GetRecordList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.GetRecordListReceiveStruct)); err == nil {
		results, err := users.GetRecordList(rec, uid)
		us.Response(ctx, results, err)
	}
}

// ClearRecord 清空历史记录
func (us UserControllers) ClearRecord(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	res, err := users.ClearRecord(uid)
	us.Response(ctx, res, err)
}

// DeleteRecordByID 根据数据库的ID来删除某一条记录
func (us UserControllers) DeleteRecordByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.DeleteRecordByIDReceiveStruct)); err == nil {
		results, err := users.DeleteRecordByID(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetNoticeList 获取通知列表
func (us UserControllers) GetNoticeList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.GetNoticeListReceiveStruct)); err == nil {
		results, err := users.GetNoticeList(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetChatList 获取聊天列表
func (us UserControllers) GetChatList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := users.GetChatList(uid)
	us.Response(ctx, results, err)
}

// GetChatHistoryMsg 获取历史聊天记录
func (us UserControllers) GetChatHistoryMsg(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.GetChatHistoryMsgStruct)); err == nil {
		results, err := users.GetChatHistoryMsg(rec, uid)
		us.Response(ctx, results, err)
	}
}

// PersonalLetter 点击私信时触发
func (us UserControllers) PersonalLetter(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.PersonalLetterReceiveStruct)); err == nil {
		results, err := users.PersonalLetter(rec, uid)
		us.Response(ctx, results, err)
	}
}

// DeleteChatItem 删除聊天记录
func (us UserControllers) DeleteChatItem(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(receive.DeleteChatItemReceiveStruct)); err == nil {
		results, err := users.DeleteChatItem(rec, uid)
		us.Response(ctx, results, err)
	}
}

// CheckIn 签到
func (us UserControllers) CheckIn(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, &receive.CheckInRequestStruct{UID: uid}); err == nil {
		results, err := checkin.CheckIn(rec)
		us.Response(ctx, results, err)
	}
}

// GetIntegral 获取用户积分
func (us UserControllers) GetIntegral(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, &receive.GetUserIntegralRequest{UID: uid}); err == nil {
		results, err := checkin.GetUserIntegral(rec)
		us.Response(ctx, results, err)
	}
}
