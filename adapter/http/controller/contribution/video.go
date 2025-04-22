package contribution

import (
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive/contribution/video"
	video2 "fakebilibili/domain/servicce/contribution/video"
	"fakebilibili/infrastructure/model/contribution/video/comments"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/response"
	"fakebilibili/infrastructure/pkg/utils/validator"
	video3 "fakebilibili/quartzImpl/video"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"time"
)

type Controllers struct {
	controller.BaseControllers
}

// GetVideoBarrage  获取视频弹幕 (播放器）
func (c Controllers) GetVideoBarrage(ctx *gin.Context) {
	GetVideoBarrageRec := new(video.GetVideoBarrageReceiveStruct)
	GetVideoBarrageRec.ID = ctx.Query("id")
	results, err := video2.GetVideoBarrage(GetVideoBarrageRec)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.BarrageSuccess(ctx, results)
}

// GetVideoBarrageList  获取视频弹幕展示（先从redis中查找）
func (c Controllers) GetVideoBarrageList(ctx *gin.Context) {
	GetVideoBarrageRec := new(video.GetVideoBarrageListReceiveStruct)
	GetVideoBarrageRec.ID = ctx.Query("id")
	results, err := video2.GetVideoBarrageList(GetVideoBarrageRec)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.BarrageSuccess(ctx, results)
}

// GetVideoComment 获取视频评论
func (c Controllers) GetVideoComment(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(video.GetVideoCommentReceiveStruct)); err == nil {
		results, err := video2.GetVideoComment(rec)
		c.Response(ctx, results, err)
	}
}

// GetVideoCommentCountById 根据视频id返回视频的评论总条数
func (c Controllers) GetVideoCommentCountById(ctx *gin.Context) {
	/*
		踩了几个雷：
		1、前端post方法传递过来的参数{id:"3123"}不能用ctx.Query("id")取，要构造一个json对象，然后ctx.BindJSON取参数id
		2、comments表的videoId字段为uint类型，取出id之后还要转为uint类型才能正确的进行查询
	*/
	var json struct {
		ID string `json:"id"`
	}
	if err := ctx.BindJSON(&json); err != nil {
		global.Logger.Errorf("类型转换错误：%v", err)
	}
	value, _ := strconv.Atoi(json.ID)
	videoId := uint(value)
	var count int64
	err := global.MysqlDb.Model(&comments.Comment{}).Where("video_id = ?", videoId).Count(&count).Error
	c.Response(ctx, count, err)
}

// LikeVideoComment 给评论点赞
func (c Controllers) LikeVideoComment(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(video.LikeVideoCommentReqStruct)); err == nil {
		results, err := video2.LikeVideoComment(rec)
		c.Response(ctx, results, err)
	}
}

// GetVideoContributionByID 根据id获取视频信息
// 每次获取视频信息，该视频热度/播放量Heat应+1
func (c Controllers) GetVideoContributionByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(video.GetVideoContributionByIDReceiveStruct)); err == nil {
		res, err := video2.GetVideoContributionByID(rec, uid)
		c.Response(ctx, res, err)
	}
}

// SendVideoBarrage 发送视频弹幕
// 1.保存弹幕到数据库 2.通过socket通知videoId为标识房间里的所有用户
func (c Controllers) SendVideoBarrage(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	SendVideoBarrageRec := new(video.SendVideoBarrageReceiveStruct)
	if err := ctx.ShouldBindBodyWith(SendVideoBarrageRec, binding.JSON); err == nil {
		res, err := video2.SendVideoBarrage(SendVideoBarrageRec, uid)
		if err != nil {
			response.Error(ctx, err.Error())
			return
		}
		response.BarrageSuccess(ctx, res)
	} else {
		validator.CheckParams(ctx, err)
	}
}

// CreateVideoContribution 发布视频和编辑视频
func (c Controllers) CreateVideoContribution(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(video.CreateVideoContributionReceiveStruct)); err == nil {
		res, err := video2.CreateVideoContribution(rec, uid)
		if err != nil {
			global.Logger.Errorf("保存视频信息失败：%v", err)
		}

		// 定时发布视频
		if rec.DateTime != "" {
			// 要指定时区，不然时间不一致
			location, _ := time.LoadLocation("Local")
			targetTime, err := time.ParseInLocation("2006-01-02 15:04:05", rec.DateTime, location)
			if err != nil {
				global.Logger.Errorf("解析时间出错，传递过来的dataTime为%v", rec.DateTime)
			}
			err = video3.PublishVideoOnSchedule(targetTime, res.(uint))
		}
		c.Response(ctx, res, err)
	}
}

// UpdateVideoContribution 更新视频信息
func (c Controllers) UpdateVideoContribution(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(video.UpdateVideoContributionReceiveStruct)); err == nil {
		res, err := video2.UpdateVideoContribution(rec, uid)
		c.Response(ctx, res, err)
	}
}

// DeleteVideoByID 通过ID删除视频
func (c Controllers) DeleteVideoByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(video.DeleteVideoByIDReceiveStruct)); err == nil {
		res, err := video2.DeleteVideoByID(rec, uid)
		c.Response(ctx, res, err)
	}
}

// VideoPostComment 进行视频评论
func (c Controllers) VideoPostComment(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(video.VideosPostCommentReceiveStruct)); err == nil {
		res, err := video2.VideoPostComment(rec, uid)
		c.Response(ctx, res, err)
	}
}

// GetVideoManagementList 创作中心获取视频稿件列表
func (c Controllers) GetVideoManagementList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(video.GetVideoManagementListReceiveStruct)); err == nil {
		res, err := video2.GetVideoManagementList(rec, uid)
		c.Response(ctx, res, err)
	}
}

// LikeVideo 给视频点赞
func (c Controllers) LikeVideo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(video.LikeVideoReceiveStruct)); err == nil {
		res, err := video2.LikeVideo(rec, uid)
		c.Response(ctx, res, err)
	}
}

// DeleteVideoByPath 通过路径删除视频
func (c Controllers) DeleteVideoByPath(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(video.DeleteVideoByPathReceiveStruct)); err == nil {
		res, err := video2.DeleteVideoByPath(rec)
		c.Response(ctx, res, err)
	}
}

// GetLastWatchTime 返回上次观看视频的进度
func (c Controllers) GetLastWatchTime(ctx *gin.Context) {
	vid, _ := strconv.ParseInt(ctx.Query("id"), 10, 64)
	uid := ctx.GetUint("uid")
	res, err := video2.GetLastWatchTime(uid, uint(vid))
	c.Response(ctx, res, err)
}

// SendWatchTime 保存视频观看进度
func (c Controllers) SendWatchTime(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(video.SendWatchTimeReqStruct)); err == nil {
		err := video2.SendWatchTime(rec, uid)
		if err != nil {
			c.Response(ctx, "保存失败", err)
		}
		c.Response(ctx, "保存成功", err)
	}
}
