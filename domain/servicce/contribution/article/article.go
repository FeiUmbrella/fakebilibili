package article

import (
	"encoding/json"
	"fakebilibili/adapter/http/receive/contribution/article"
	article3 "fakebilibili/adapter/http/response/contribution/article"
	"fakebilibili/domain/servicce/users"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/common"
	article2 "fakebilibili/infrastructure/model/contribution/article"
	"fakebilibili/infrastructure/model/user/notice"
	"fakebilibili/infrastructure/model/user/record"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"fmt"
	"github.com/FeiUmbrella/Sensitive_Words_Filter/filter"
	"github.com/dlclark/regexp2"
	"strconv"
	"strings"
)

// GetArticleContributionList 首页显示的文章
func GetArticleContributionList(data *article.GetArticleContributionListReceiveStruct) (results interface{}, err error) {
	articlesList := &article2.ArticlesContributionList{}
	if articlesList.GetList(data.PageInfo) {
		return nil, fmt.Errorf("查询失败")
	}
	return article3.GetArticleContributionListResponse(articlesList), nil
}

// GetArticleContributionListByUser 获取用户发布的文章
func GetArticleContributionListByUser(data *article.GetArticleContributionListByUserReceiveStruct) (results interface{}, err error) {
	articlesList := &article2.ArticlesContributionList{}
	if articlesList.GetListByUid(data.UserID) {
		return nil, fmt.Errorf("查询失败")
	}
	return article3.GetArticleContributionListByUserResponse(articlesList), nil
}

// GetArticleComment 获取文章评论
func GetArticleComment(data *article.GetArticleCommentReceiveStruct) (results interface{}, err error) {
	articlesList := &article2.ArticlesContribution{}
	if articlesList.GetArticleComments(data.ArticleID, data.PageInfo) {
		return nil, fmt.Errorf("查询失败")
	}
	return article3.GetArticleContributionCommentsResponse(articlesList), nil
}

// GetArticleClassificationList 按照分类获取文章列表
func GetArticleClassificationList() (results interface{}, err error) {
	cf := &article2.ClassificationsList{}
	err = cf.FindAll()
	if err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	return article3.GetArticleClassificationListResponse(cf), nil
}

// GetArticleTotalInfo 获取文章相关总和信息
func GetArticleTotalInfo() (results interface{}, err error) {
	// 查询文章数量
	articleNum := new(int64)
	articlesList := &article2.ArticlesContributionList{}
	articlesList.GetAllCount(articleNum)

	// 查询文章分类信息
	cf := make(article2.ClassificationsList, 0)
	err = cf.FindAll()
	if err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	cfNum := int64(len(cf))
	return article3.GetArticleTotalInfoResponse(&cf, articleNum, cfNum), nil
}

// GetArticleContributionByID 根据文章id获取文章即观看某个文章
func GetArticleContributionByID(data *article.GetArticleContributionByIDReceiveStruct, uid uint) (results interface{}, err error) {
	articles := &article2.ArticlesContribution{}
	if !articles.GetInfoByID(data.ArticleID) {
		return nil, fmt.Errorf("查询失败")
	}
	if uid > 0 {
		rd := new(record.Record)
		err = rd.AddArticleRecord(uid, data.ArticleID)
		if err != nil {
			return nil, fmt.Errorf("添加观看文章历史记录失败")
		}
		// 文章热度++
		if !global.RedisDb.SIsMember(consts.ArticleWatchByID+strconv.Itoa(int(data.ArticleID)), uid).Val() {
			// uid最近没看过这篇文章
			global.RedisDb.SAdd(consts.ArticleWatchByID+strconv.Itoa(int(data.ArticleID)), uid)
			if err := articles.Watch(data.ArticleID); err != nil {
				global.Logger.Error("添加热度错误article_id:", err)
			}
			articles.Heat++
		}
	}
	return
}

// CreateArticleContribution 发布文章
func CreateArticleContribution(data *article.CreateArticleContributionReceiveStruct, uid uint) (results interface{}, err error) {
	for _, v := range data.Label {
		vRune := []rune(v)
		if len(vRune) > 7 {
			return nil, fmt.Errorf("标签长度不能大于7位")
		}
	}

	// 文章封面
	coverImg, _ := json.Marshal(common.Img{
		Tp:  data.CoverUploadType,
		Src: data.Cover,
	})

	// 正则匹配替换url
	// 取url前缀
	prefix, err := conversion.SwitchTypeAsUrlPrefix(data.ArticleContributionUploadType)
	if err != nil {
		return nil, fmt.Errorf("保存资源方式不存在")
	}

	// 正则匹配替换
	reg := regexp2.MustCompile(`(?<=(img[^>]*src="))[^"]*?`+prefix, 0)
	match, err := reg.Replace(data.Content, consts.UrlPrefixSubstitution, -1, -1)
	data.Content = match
	// 向数据库插入文章信息
	articles := &article2.ArticlesContribution{
		Uid:                uid,
		ClassificationID:   data.ClassificationID,
		Title:              data.Title,
		Cover:              coverImg,
		Label:              strings.Join(data.Label, ","),
		Content:            data.Content,
		ContentStorageType: data.ArticleContributionUploadType,
		IsComments:         conversion.Bool2Int8(*data.Comments),
		Heat:               0,
	}
	if !articles.Create() {
		return nil, fmt.Errorf("保存失败")
	}
	return "文章信息保存成功", nil
}

// UpdateArticleContribution 更新文章信息
func UpdateArticleContribution(data *article.UpdateArticleContributionReceiveStruct, uid uint) (results interface{}, err error) {
	articles := &article2.ArticlesContribution{}
	if !articles.GetInfoByID(data.ID) {
		return nil, fmt.Errorf("更新视频不存在")
	}
	if articles.Uid != uid {
		return nil, fmt.Errorf("非法操作")
	}
	coverImg, _ := json.Marshal(common.Img{
		Tp:  data.CoverUploadType,
		Src: data.Cover,
	})

	updateList := map[string]interface{}{
		"cover":             coverImg,
		"title":             data.Title,
		"content":           data.Content,
		"label":             strings.Join(data.Label, ","),
		"is_comments":       data.Comments,
		"classification_id": data.ClassificationID,
	}
	// 进行视频资料更新
	if !articles.Update(updateList) {
		return nil, fmt.Errorf("更新数据失败")
	}
	return "文章信息更新成功", nil
}

// DeleteArticleByID 删除文章
func DeleteArticleByID(data *article.DeleteArticleByIDReceiveStruct, uid uint) (results interface{}, err error) {
	ac := new(article2.ArticlesContribution)
	if !ac.Delete(data.ID, uid) {
		return nil, fmt.Errorf("删除失败")
	}
	return "删除成功", nil
}

// ArticlePostComment 发表文章评论
func ArticlePostComment(data *article.ArticlesPostCommentReceiveStruct, uid uint) (results interface{}, err error) {
	node := filter.NewNode()
	node.StartMatch(data.Content)
	if node.IsSensitive() {
		return nil, fmt.Errorf("评论中包含敏感词")
	}
	articleInfo := new(article2.ArticlesContribution)
	if !articleInfo.GetInfoByID(data.ArticleID) {
		return nil, fmt.Errorf("评论文章不存在")
	}

	// 找到被评论的那条评论
	ct := article2.Comment{}
	ct.ID = data.ContentID

	// 获取被评论那条评论的根评论
	CommentFirstID := ct.GetCommentFirstID()

	// 找到被评论那条评论的所属用户
	ctu := article2.Comment{}
	ctu.ID = data.ContentID
	commentUserID := ctu.GetCommentUserID()

	comment := article2.Comment{
		Uid:            uid,
		ArticleID:      data.ArticleID,
		Context:        data.Content,
		CommentID:      data.ContentID, // 父评论
		CommentUserID:  commentUserID,  // 父评论的用户
		CommentFirstID: CommentFirstID, // 根评论
	}
	if !comment.Create() {
		return nil, fmt.Errorf("发布文章评论失败")
	}

	// socket 推送(文章发布者在线情况下)
	if _, ok := users.Severe.UserMapChannel[articleInfo.UserInfo.ID]; ok {
		userChannel := users.Severe.UserMapChannel[articleInfo.UserInfo.ID]
		userChannel.NoticeMessage(notice.ArticleComment)
	}
	return "发布文章评论成功", nil
}

// GetArticleManagementList 创作中心获取专栏稿件列表
func GetArticleManagementList(data *article.GetArticleManagementListReceiveStruct, uid uint) (results interface{}, err error) {
	//获取个人发布专栏信息
	list := new(article2.ArticlesContributionList)
	err = list.GetArticleManagementList(data.PageInfo, uid)
	if err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	res, err := article3.GetArticleManagementListResponse(list)
	if err != nil {
		return nil, fmt.Errorf("响应失败")
	}
	return res, nil
}

// GetColumnByClassificationId 根据专栏id获取用户创建的对应分类下的专栏
func GetColumnByClassificationId(data *article.GetColumnByClassificationId, uid uint) (results interface{}, err error) {
	acl := new(article2.ArticlesContributionList)
	err = global.MysqlDb.
		Where("uid = ? AND classification_id = ?", uid, data.ClassificationID).
		Preload("Likes").
		Preload("Comments").
		Preload("Classification").
		Order("created_at desc").
		Find(&acl).Error
	if err != nil {
		return nil, fmt.Errorf("根据文章类别ID查询文章失败")
	}
	return article3.GetArticleContributionListResponse(acl), nil
}
