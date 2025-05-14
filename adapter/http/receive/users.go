package receive

import (
	"fakebilibili/infrastructure/model/common"
	"time"
)

// RegisterReceiveStruct 用户注册相关参数结构体
type RegisterReceiveStruct struct {
	UserName         string `json:"username" binding:"required" form:"username"`
	Password         string `json:"password" binding:"required" form:"password"`
	Email            string `json:"email" binding:"required,email" form:"email"`
	VerificationCode string `json:"verificationCode" binding:"required" form:"verificationCode"` // 邮箱验证码
}

// SendEmailVerCodeReceiveStruct 发送邮箱验证码
type SendEmailVerCodeReceiveStruct struct {
	Email string `json:"email" binding:"required,email" form:"email"`
}

// LoginReceiveStruct 用户登录 相关参数结构体
type LoginReceiveStruct struct {
	UserName string `json:"username" binding:"required" form:"username"`
	Password string `json:"password" binding:"required" form:"password"`
}

// ForgetReceiveStruct 忘记密码进行找回 相关参数结构体
type ForgetReceiveStruct struct {
	Password         string `json:"password" binding:"required" form:"password"`
	Email            string `json:"email" binding:"required,email" form:"email"`
	VerificationCode string `json:"verificationCode" binding:"required" form:"verificationCode"` // 邮箱验证码
}

// GetAttentionListReceiveStruct 获取用户关注列表 相关参数结构体
type GetAttentionListReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// GetVermicelliListReceiveStruct 获取用户粉丝列表 相关参数结构体
type GetVermicelliListReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// GetSpaceIndividualReceiveStruct 获取用户个人空间 相关参数结构体
type GetSpaceIndividualReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// GetReleaseInformationReceiveStruct 获取用户发布的作品(视频and专栏) 相关参数结构体
type GetReleaseInformationReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// SetUserInfoReceiveStruct 接收设置用户信息参数 结构体
type SetUserInfoReceiveStruct struct {
	Username  string    `json:"username" binding:"required" form:"username"`
	Gender    *int      `json:"gender" binding:"required" form:"gender"`
	BirthDate time.Time `json:"birth_Date" binding:"required" form:"birth_date"`
	IsVisible *bool     `json:"is_Visible" binding:"required" form:"is_visible"`
	//个性签名和社媒地址不应该设置为required，没填的时候会绑定失败
	Signature   string `json:"signature" form:"signature"`
	SocialMedia string `json:"social_media" form:"social_media"`
}

// DetermineNameExistsStruct 改名时传进来的名字 结构体
type DetermineNameExistsStruct struct {
	Username string `json:"username" binding:"required" form:"username"`
}

// UpdateAvatarStruct 修改头像 结构体
type UpdateAvatarStruct struct {
	ImgUrl string `json:"imgUrl" binding:"required" form:"imgUrl"`
	Tp     string `json:"type" binding:"required" form:"type"`
}

// SaveLiveDataReceiveStruct 修改直播间信息 结构体
type SaveLiveDataReceiveStruct struct {
	Tp     string `json:"type" binding:"required" form:"type"`
	ImgUrl string `json:"imgUrl" binding:"required" form:"img_url"`
	Title  string `json:"title" binding:"required" form:"title"`
}

// ChangePasswordReceiveStruct 登录状态下修改密码
type ChangePasswordReceiveStruct struct {
	VerificationCode string `json:"verificationCode" binding:"required" form:"verificationCode"`
	Password         string `json:"password" binding:"required" form:"password"`
	ConfirmPassword  string `json:"confirm_password" binding:"required" form:"confirm_password"`
}

// AttentionReceiveStruct 关注用户时参数 结构体
type AttentionReceiveStruct struct {
	Uid uint `json:"uid"  binding:"required" binding:"required" form:"uid"`
}

// CreateFavoritesReceiveStruct 创建收藏夹所需参数
type CreateFavoritesReceiveStruct struct {
	ID      uint   `json:"id" form:"id"`
	Title   string `json:"title" binding:"required" form:"title"`
	Content string `json:"content" form:"content"`
	Cover   string `json:"cover" form:"cover"`
	Tp      string `json:"type" form:"type"`
}

// DeleteFavoritesReceiveStruct 删除收藏夹所需参数
type DeleteFavoritesReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// FavoriteVideoReceiveStruct 收藏/取消收藏视频
// 设VideoID原来所在的收藏夹的id集合为old_ids，如果IDs中id不在old_ids中，那么是要收藏该视频到为id的收藏夹中
// 如果old_ids中id不在IDs中，那么是要取消收藏在id的收藏夹中的该视频
type FavoriteVideoReceiveStruct struct {
	IDs     []uint `json:"ids" binding:"required" form:"ids"`           //视频应该位于的收藏夹的ids
	VideoID uint   `json:"video_id" binding:"required" form:"video_id"` //视频的video_id
}

// GetFavoritesListByFavoriteVideoReceiveStruct 获取用户包含某个视频的收藏夹列表
type GetFavoritesListByFavoriteVideoReceiveStruct struct {
	VideoID uint `json:"video_id" binding:"required" form:"video_id"`
}

// GetFavoriteVideoListReceiveStruct 获取收藏夹视频列表
type GetFavoriteVideoListReceiveStruct struct {
	FavoriteID uint `json:"favorite_id" binding:"required" form:"favorite_id"`
}

// GetCollectListNameReceiveStruct 获取某个收藏夹名字
type GetCollectListNameReceiveStruct struct {
	FavoriteID uint `json:"favorite_id" binding:"required" form:"favorite_id"`
}

// GetRecordListReceiveStruct 获取历史记录
type GetRecordListReceiveStruct struct {
	PageInfo common.PageInfo `json:"page_info" binding:"required" form:"page_info"`
}

// DeleteRecordByIDReceiveStruct 根据ID删除一条历史记录
type DeleteRecordByIDReceiveStruct struct {
	ID uint `json:"id" form:"id"`
}

// GetNoticeListReceiveStruct 获取通知列表
type GetNoticeListReceiveStruct struct {
	Type     string          `json:"type" form:"type"`
	PageInfo common.PageInfo `json:"page_info" binding:"required" form:"page_info"`
}

// GetChatHistoryMsgStruct 获取聊天记录
type GetChatHistoryMsgStruct struct {
	Tid      uint      `json:"tid" form:"tid"`
	LastTime time.Time `json:"last_time" form:"last_time"`
}

// PersonalLetterReceiveStruct 私信
type PersonalLetterReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// DeleteChatItemReceiveStruct 删除聊天记录
type DeleteChatItemReceiveStruct struct {
	ID uint `json:"id" binding:"required" form:"id"`
}

// CheckInRequestStruct 签到
type CheckInRequestStruct struct {
	UID uint `json:"uid" binding:"required"`
}

// GetUserIntegralRequest 获取用户积分
type GetUserIntegralRequest struct {
	UID uint `json:"uid" binding:"required"`
}

// GetCaptchaStruct 请求验证码携带的captchaId
type GetCaptchaStruct struct {
	CaptchaId string `json:"captcha_id" binding:"required"`
}
