package users

import (
	"fakebilibili/adapter/http/controller/users"
	"fakebilibili/adapter/http/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *LoginRouter) InitRouter(Router *gin.RouterGroup) {
	router := Router.Group("user").Use(middleware.VerificationToken())
	{
		userControllers := new(users.UserControllers)
		router.POST("/getUserInfo", userControllers.GetUserInfo)
		router.POST("/setUserInfo", userControllers.SetUserInfo)
		router.POST("/determineNameExists", userControllers.DetermineNameExists)
		router.POST("/updateAvatar", userControllers.UpdateAvatar)
		router.POST("/getLiveData", userControllers.GetLiveInfo)
		router.POST("/saveLiveData", userControllers.SaveLiveInfo)
		router.POST("/sendEmailVerificationCodeByChangePassword", userControllers.SendEmailVerificationCodeByChangePassword)
		router.POST("/changePassword", userControllers.ChangePassword)
		router.POST("/attention", userControllers.Attention)
		router.POST("/createFavorites", userControllers.CreateFavorites)
		router.POST("/getFavoritesList", userControllers.GetFavoritesList)
		router.POST("/deleteFavorites", userControllers.DeleteFavorites)
		router.POST("/favoriteVideo", userControllers.FavoriteVideo)
		router.POST("/getFavoritesListByFavoriteVideo", userControllers.GetFavoritesListByFavoriteVideo)
		router.POST("/getFavoriteVideoList", userControllers.GetFavoriteVideoList)
		router.POST("/user/getCollectListName", userControllers.GetCollectListName)
		router.POST("/getRecordList", userControllers.GetRecordList)
		router.POST("/clearRecord", userControllers.ClearRecord)
		router.POST("/deleteRecordByID", userControllers.DeleteRecordByID)
		router.POST("/getNoticeList", userControllers.GetNoticeList)
		router.POST("/getChatList", userControllers.GetChatList)
		router.POST("/getChatHistoryMsg", userControllers.GetChatHistoryMsg)
		router.POST("/personalLetter", userControllers.PersonalLetter)
		router.POST("/deleteChatItem", userControllers.DeleteChatItem)
		router.POST("/checkin", userControllers.CheckIn)         //用户签到
		router.POST("/getIntegral", userControllers.GetIntegral) //获取用户的积分
	}
}
