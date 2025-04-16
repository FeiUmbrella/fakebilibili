package osscommonality

import (
	"fakebilibili/adapter/http/controller"
	osscommonality2 "fakebilibili/adapter/http/receive/osscommonality"
	"fakebilibili/domain/servicce/osscommonality"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	controller.BaseControllers
}

// OssSTS 利用阿里云OSS 的 Security Token Service(STS)获取临时授权token 来上传文件到OSS
func (c *Controllers) OssSTS(ctx *gin.Context) {
	res, err := osscommonality.OssSTS()
	c.Response(ctx, res, err)
}

// Upload 文件上传
func (c *Controllers) Upload(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	res, err := osscommonality.Upload(file, ctx)
	c.Response(ctx, res, err)
}

// UploadSlice  上传文件的一个分片
func (c *Controllers) UploadSlice(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	results, err := osscommonality.UploadSlice(file, ctx)
	c.Response(ctx, results, err)
}

// UploadCheck 验证文件是否已经上传，若为上传返回文件的哪些分片未上传
func (c *Controllers) UploadCheck(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(osscommonality2.UploadCheckStruct)); err == nil {
		results, err := osscommonality.UploadCheck(rec)
		c.Response(ctx, results, err)
	}
}

// UploadMerge 合并分片上传保存
func (c *Controllers) UploadMerge(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(osscommonality2.UploadMergeStruct)); err == nil {
		results, err := osscommonality.UploadMerge(rec)
		c.Response(ctx, results, err)
	}
}

// UploadingMethod 获取上传文件配置（上传目录、图片文件的应该保存的大小）
func (c *Controllers) UploadingMethod(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(osscommonality2.UploadingMethodStruct)); err == nil {
		results, err := osscommonality.UploadingMethod(rec)
		c.Response(ctx, results, err)
	}
}

// UploadingDir 获取上传文件保存目录
func (c *Controllers) UploadingDir(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(osscommonality2.UploadingDirStruct)); err == nil {
		results, err := osscommonality.UploadingDir(rec)
		c.Response(ctx, results, err)
	}
}

// GetFullPathOfImage 获取图片完整路径
func (c *Controllers) GetFullPathOfImage(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(osscommonality2.GetFullPathOfImageMethodStruct)); err == nil {
		results, err := osscommonality.GetFullPathOfImage(rec)
		c.Response(ctx, results, err)
	}
}

// Search 用关键词搜索用户/视频
func (c *Controllers) Search(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	//SearchStruct
	if rec, err := controller.ShouldBind(ctx, new(osscommonality2.SearchStruct)); err == nil {
		results, err := osscommonality.Search(rec, uid)
		c.Response(ctx, results, err)
	}
}

// RegisterMedia 注册媒体资源,使用阿里云的只能媒体服务
func (c *Controllers) RegisterMedia(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(osscommonality2.RegisterMediaStruct)); err == nil {
		results, err := osscommonality.RegisterMedia(rec)
		c.Response(ctx, results, err)
	}
}
