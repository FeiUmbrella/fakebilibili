package contribution

import (
	"fakebilibili/adapter/http/controller"
	"fakebilibili/adapter/http/receive/contribution/article"
	article2 "fakebilibili/domain/servicce/contribution/article"
	"github.com/gin-gonic/gin"
)

// GetArticleContributionList 首页查询专栏
func (c Controllers) GetArticleContributionList(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(article.GetArticleContributionListReceiveStruct)); err == nil {
		res, err := article2.GetArticleContributionList(rec)
		c.Response(ctx, res, err)
	}
}

// GetArticleContributionListByUser 获取用户的文章
func (c Controllers) GetArticleContributionListByUser(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(article.GetArticleContributionListByUserReceiveStruct)); err == nil {
		res, err := article2.GetArticleContributionListByUser(rec)
		c.Response(ctx, res, err)
	}
}

// GetArticleComment 获取文章评论
func (c Controllers) GetArticleComment(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(article.GetArticleCommentReceiveStruct)); err == nil {
		res, err := article2.GetArticleComment(rec)
		c.Response(ctx, res, err)
	}
}

// GetArticleClassificationList 按照分类获取文章列表
func (c Controllers) GetArticleClassificationList(ctx *gin.Context) {
	res, err := article2.GetArticleClassificationList()
	c.Response(ctx, res, err)
}

// GetArticleTotalInfo 获取文章相关总和信息
func (c Controllers) GetArticleTotalInfo(ctx *gin.Context) {
	res, err := article2.GetArticleTotalInfo()
	c.Response(ctx, res, err)
}

// GetArticleContributionByID 根据文章id获取文章根据文章id获取文章即观看某个文章
func (c Controllers) GetArticleContributionByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(article.GetArticleContributionByIDReceiveStruct)); err == nil {
		res, err := article2.GetArticleContributionByID(rec, uid)
		c.Response(ctx, res, err)
	}
}

// CreateArticleContribution 发布视频
func (c Controllers) CreateArticleContribution(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(article.CreateArticleContributionReceiveStruct)); err == nil {
		res, err := article2.CreateArticleContribution(rec, uid)
		c.Response(ctx, res, err)
	}
}

// UpdateArticleContribution 更新文章信息
func (c Controllers) UpdateArticleContribution(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(article.UpdateArticleContributionReceiveStruct)); err == nil {
		res, err := article2.UpdateArticleContribution(rec, uid)
		c.Response(ctx, res, err)
	}
}

// DeleteArticleByID 删除专栏
func (c Controllers) DeleteArticleByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(article.DeleteArticleByIDReceiveStruct)); err == nil {
		res, err := article2.DeleteArticleByID(rec, uid)
		c.Response(ctx, res, err)
	}
}

// ArticlePostComment 发布文章评论
func (c Controllers) ArticlePostComment(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(article.ArticlesPostCommentReceiveStruct)); err == nil {
		res, err := article2.ArticlePostComment(rec, uid)
		c.Response(ctx, res, err)
	}
}

// GetArticleManagementList 创作中心获取专栏稿件列表
func (c Controllers) GetArticleManagementList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(article.GetArticleManagementListReceiveStruct)); err == nil {
		results, err := article2.GetArticleManagementList(rec, uid)
		c.Response(ctx, results, err)
	}
}

// GetColumnByClassificationId 根据专栏id获取用户创建的对应分类下的专栏
func (c Controllers) GetColumnByClassificationId(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controller.ShouldBind(ctx, new(article.GetColumnByClassificationId)); err == nil {
		results, err := article2.GetColumnByClassificationId(rec, uid)
		c.Response(ctx, results, err)
	}
}
