package response

import (
	"fakebilibili/infrastructure/model/contribution/article"
	"fakebilibili/infrastructure/model/contribution/video"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/attention"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"github.com/dlclark/regexp2"
	"time"
)

// UserInfoResponseStruct 返回前端用户信息结构体
type UserInfoResponseStruct struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"username"`
	Photo     string    `json:"photo"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

// UserInfoResponse  生成返回用用户信息结构体
func UserInfoResponse(us *user.User, token string) UserInfoResponseStruct {
	//判断用户是否为微信用户进行图片处理
	photo, _ := conversion.FormattingJsonSrc(us.Photo)
	return UserInfoResponseStruct{
		ID:        us.ID,
		UserName:  us.Username,
		Photo:     photo,
		Token:     token,
		CreatedAt: us.CreatedAt,
	}
}

// GetAttentionListInfo 返回前端的关注列表结构体
type GetAttentionListInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Signature   string `json:"signature"`
	Photo       string `json:"photo"`
	IsAttention bool   `json:"is_attention"`
}
type GetAttentionListInfoList []GetAttentionListInfo

// GetAttentionListResponse 生成关注列表结构体并返回
func GetAttentionListResponse(atl *attention.AttentionsList, arr []uint) (interface{}, error) {
	list := make(GetAttentionListInfoList, 0)
	for _, at := range *atl {
		photo, _ := conversion.FormattingJsonSrc(at.AttentionUserInfo.Photo)
		isAttention := false
		for _, ak := range arr {
			if ak == at.AttentionID {
				isAttention = true
			}
		}
		list = append(list, GetAttentionListInfo{
			ID:          at.AttentionID,
			Name:        at.AttentionUserInfo.Username,
			Signature:   at.AttentionUserInfo.Signature,
			Photo:       photo,
			IsAttention: isAttention,
		})
	}
	return list, nil
}

// GetVermicelliListInfo 返回前端的粉丝列表结构体
type GetVermicelliListInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Signature   string `json:"signature"`
	Photo       string `json:"photo"`
	IsAttention bool   `json:"is_attention"`
}

type GetVermicelliListInfoList []GetVermicelliListInfo

// GetVermicelliListResponse 生成粉丝列表结构体并返回
func GetVermicelliListResponse(al *attention.AttentionsList, arr []uint) (data interface{}, err error) {
	list := make(GetVermicelliListInfoList, 0)
	for _, v := range *al {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo)
		isAttention := false
		for _, ak := range arr {
			if ak == v.Uid {
				isAttention = true
			}
		}
		list = append(list, GetVermicelliListInfo{
			ID:          v.Uid,
			Name:        v.UserInfo.Username,
			Signature:   v.UserInfo.Signature,
			Photo:       photo,
			IsAttention: isAttention,
		})
	}
	return list, nil
}

// GetSpaceIndividualResponseStruct 返回个人空间信息结构体
type GetSpaceIndividualResponseStruct struct {
	ID            uint   `json:"id"`
	UserName      string `json:"username"`
	Photo         string `json:"photo"`
	Signature     string `json:"signature"`
	IsAttention   bool   `json:"is_attention"`
	AttentionNum  *int64 `json:"attention_num"`
	VermicelliNum *int64 `json:"vermicelli_num"`
}

// GetSpaceIndividualResponse 返回个人空间信息
func GetSpaceIndividualResponse(us *user.User, isAttention bool, atNum, vmNum int64) (data interface{}, err error) {
	photo, _ := conversion.FormattingJsonSrc(us.Photo)
	return GetSpaceIndividualResponseStruct{
		ID:            us.ID,
		UserName:      us.Username,
		Photo:         photo,
		Signature:     us.Signature,
		IsAttention:   isAttention,
		AttentionNum:  &atNum,
		VermicelliNum: &vmNum,
	}, nil
}

// ReleaseInformationVideoInfo 发布的视频信息 结构体
type ReleaseInformationVideoInfo struct {
	ID            uint      `json:"id"`
	Uid           uint      `json:"uid" `
	Title         string    `json:"title" `
	Video         string    `json:"video"`
	Cover         string    `json:"cover" `
	VideoDuration int64     `json:"video_duration"`
	Label         []string  `json:"label"`
	Introduce     string    `json:"introduce"`
	Heat          int       `json:"heat"`
	BarrageNumber int       `json:"barrageNumber"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
}
type ReleaseInformationVideoInfoList []ReleaseInformationVideoInfo

// ReleaseInformationArticleInfo 发布的专栏信息 结构体
type ReleaseInformationArticleInfo struct {
	Id             uint      `json:"id"`
	Uid            uint      `json:"uid" `
	Title          string    `json:"title" `
	Cover          string    `json:"cover" `
	Label          []string  `json:"label" `
	Content        string    `json:"content"`
	IsComments     bool      `json:"is_comments"`
	Heat           int       `json:"heat"`
	LikesNumber    int       `json:"likes_number"`
	CommentsNumber int       `json:"comments_number"`
	Classification string    `json:"classification"`
	CreatedAt      time.Time `json:"created_at"`
}
type ReleaseInformationArticleInfoList []ReleaseInformationArticleInfo

type GetReleaseInformationResponseStruct struct {
	VideoList   ReleaseInformationVideoInfoList   `json:"videoList"`
	ArticleList ReleaseInformationArticleInfoList `json:"articleList"`
}

// GetReleaseInformationResponse 返回用户发布的视频和专栏信息
func GetReleaseInformationResponse(videoList *video.VideosContributionList, articleList *article.ArticlesContributionList) (interface{}, error) {
	// 处理视频列表
	vList := make(ReleaseInformationVideoInfoList, 0)
	for _, v := range *videoList {
		cover, _ := conversion.FormattingJsonSrc(v.Cover)
		videoScr, _ := conversion.FormattingJsonSrc(v.Video)
		vList = append(vList, ReleaseInformationVideoInfo{
			ID:            v.ID,
			Uid:           v.Uid,
			Title:         v.Title,
			Cover:         cover,
			Video:         videoScr,
			VideoDuration: v.VideoDuration,
			Label:         conversion.StringConversionMap(v.Label),
			Introduce:     v.Introduce,
			Heat:          v.Heat,
			BarrageNumber: len(v.Barrage),
			Username:      v.UserInfo.Username,
			CreatedAt:     v.CreatedAt,
		})
	}

	// 处理专栏列表
	aL := make(ReleaseInformationArticleInfoList, 0)
	for _, a := range *articleList {
		coverSrc, _ := conversion.FormattingJsonSrc(a.Cover)

		//正则替换第一行中的所有的html标签符
		// <div class="title">Hello, <b>世界</b>！</div> --> Hello, 世界
		reg := regexp2.MustCompile(`<(\S*?)[^>]*>.*?|<.*? />`, 0)
		// 从头开始匹配所有的标签替换为空字符
		match, _ := reg.Replace(a.Content, "", -1, -1)
		// 如果文档内容过长，只保留前100字
		matchRune := []rune(match)
		if len(matchRune) > 100 {
			a.Content = string(matchRune[:100]) + "..."
		} else {
			a.Content = match
		}

		// 只保留2个标签
		label := conversion.StringConversionMap(a.Label)
		if len(label) > 2 {
			label = label[:2]
		}
		aL = append(aL, ReleaseInformationArticleInfo{
			Id:             a.ID,
			Uid:            a.Uid,
			Title:          a.Title,
			Cover:          coverSrc,
			Label:          label,
			Content:        a.Content,
			Classification: a.Classification.Label,
			Heat:           a.Heat,
			LikesNumber:    len(a.Likes),
			CommentsNumber: len(a.Comments),
			CreatedAt:      a.CreatedAt,
			IsComments:     a.IsComments > 0,
		})
	}
	return GetReleaseInformationResponseStruct{
		vList,
		aL,
	}, nil
}
