package response

import (
	"encoding/json"
	"fakebilibili/infrastructure/model/common"
	"fakebilibili/infrastructure/model/contribution/article"
	"fakebilibili/infrastructure/model/contribution/video"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/attention"
	"fakebilibili/infrastructure/model/user/chat"
	"fakebilibili/infrastructure/model/user/collect"
	"fakebilibili/infrastructure/model/user/favorites"
	"fakebilibili/infrastructure/model/user/liveInfo"
	"fakebilibili/infrastructure/model/user/notice"
	"fakebilibili/infrastructure/model/user/record"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"fmt"
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

// GetUserInfoResponseStruct 返回用户个人信息(性别、个性签名...) 结构体
type GetUserInfoResponseStruct struct {
	ID          uint      `json:"id"`
	UserName    string    `json:"username"`
	Gender      int8      `json:"gender"`
	BirthDate   time.Time `json:"birth_date"`
	IsVisible   bool      `json:"is_visible"`
	Signature   string    `json:"signature"`
	Email       string    `json:"email"`
	SocialMedia string    `json:"social_media"`
}

func GetUserInfoResponse(us *user.User) GetUserInfoResponseStruct {
	return GetUserInfoResponseStruct{
		ID:          us.ID,
		UserName:    us.Username,
		Gender:      us.Gender,
		BirthDate:   us.BirthDate,
		IsVisible:   us.IsVisible > 0,
		Signature:   us.Signature,
		Email:       us.Email,
		SocialMedia: us.SocialMedia,
	}
}

// GetLiveInfoResponseStruct 返回直播间信息结构体
type GetLiveInfoResponseStruct struct {
	Img   string `json:"img"`
	Title string `json:"title"`
}

func GetLiveInfoResponse(lI *liveInfo.LiveInfo) (any, error) {
	src, err := conversion.FormattingJsonSrc(lI.Img)
	if err != nil {
		return nil, fmt.Errorf("返回直播间信息时 JSON Format Error")
	}
	return GetLiveInfoResponseStruct{
		Img:   src,
		Title: lI.Title,
	}, nil
}

// GetFavoritesInfo 返回收藏夹信息
type GetFavoritesInfo struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Cover    string `json:"cover"`
	Tp       string `json:"type"`
	Src      string `json:"src"`
	Max      int    `json:"max"`
	UsesInfo struct {
		Username string `json:"username"`
	} `json:"userInfo"`
}
type GetFavoritesInfoList []GetFavoritesInfo

// GetFavoritesListResponse 返回所有收藏夹
func GetFavoritesListResponse(fl *favorites.FavoriteList) (interface{}, error) {
	list := make(GetFavoritesInfoList, 0)
	for _, fs := range *fl {
		coverInfo := new(common.Img)

		_ = json.Unmarshal(fs.Cover, coverInfo)
		cover, _ := conversion.FormattingJsonSrc(fs.Cover)
		list = append(list, GetFavoritesInfo{
			ID:      fs.ID,
			Title:   fs.Title,
			Cover:   cover,
			Content: fs.Content,
			Tp:      coverInfo.Tp,
			Src:     coverInfo.Src,
			Max:     fs.Max,
			UsesInfo: struct {
				Username string `json:"username"`
			}{Username: fs.UserInfo.Username},
		})
	}
	return list, nil
}

type GetFavoritesListByFavoriteVideoInfo struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Cover    string `json:"cover"`
	Tp       string `json:"type"`
	Src      string `json:"src"`
	Max      int    `json:"max"`
	Selected bool   `json:"selected"`
	Present  int    `json:"present"`
	UserInfo struct {
		Username string `json:"username"`
	} `json:"userInfo"`
}
type GetFavoritesListByFavoriteVideoInfoList []GetFavoritesListByFavoriteVideoInfo

// GetFavoritesListByFavoriteVideoResponse 获取用户包含某个视频的收藏夹列表
func GetFavoritesListByFavoriteVideoResponse(fl *favorites.FavoriteList, fids []uint) (interface{}, error) {
	list := make(GetFavoritesListByFavoriteVideoInfoList, 0)
	// 枚举用户的所有收藏夹
	for _, f := range *fl {
		coverInfo := new(common.Img)
		_ = json.Unmarshal(f.Cover, coverInfo)
		cover, _ := conversion.FormattingJsonSrc(f.Cover) // 收藏夹封面
		// 判断是否已选该收藏夹是否在包含视频的收藏夹fid列表中
		selected := false
		for _, fid := range fids {
			if fid == f.ID {
				selected = true
			}
		}

		list = append(list, GetFavoritesListByFavoriteVideoInfo{
			ID:       f.ID,
			Title:    f.Title,
			Content:  f.Content,
			Cover:    cover,
			Tp:       coverInfo.Tp,
			Src:      coverInfo.Src,
			Max:      f.Max,
			Selected: selected,           // 该收藏夹不包含目标视频
			Present:  len(f.CollectList), // 该收藏夹的长度-收藏视频个数
			UserInfo: struct {
				Username string `json:"username"`
			}{Username: f.UserInfo.Username},
		})
	}
	return list, nil
}

type GetFavoriteVideoListItem struct {
	ID            uint      `json:"id"`
	Uid           uint      `json:"uid"`
	Title         string    `json:"title"`
	Video         string    `json:"video"`
	Cover         string    `json:"cover"`
	VideoDuration int64     `json:"video_duration"`
	CreatedAt     time.Time `json:"created_at"`
}
type GetFavoriteVideoList []GetFavoriteVideoListItem
type GetFavoriteVideoListResponseStruct struct {
	VideoList GetFavoriteVideoList `json:"videoList"`
}

// GetFavoriteVideoListResponse 获取收藏夹中的视频列表
func GetFavoriteVideoListResponse(cl *collect.CollectsList) (data interface{}, err error) {
	vl := make(GetFavoriteVideoList, 0)
	for _, c := range *cl {
		collectVideo := c.VideoInfo
		videoCover, _ := conversion.FormattingJsonSrc(collectVideo.Cover)
		videoSrc, _ := conversion.FormattingJsonSrc(collectVideo.Video)
		vl = append(vl, GetFavoriteVideoListItem{
			ID:            collectVideo.ID,
			Uid:           collectVideo.Uid,
			Title:         collectVideo.Title,
			VideoDuration: collectVideo.VideoDuration,
			Video:         videoSrc,
			Cover:         videoCover,
			CreatedAt:     c.CreatedAt,
		})
	}
	return GetFavoriteVideoListResponseStruct{vl}, nil
}

type GetRecordListItem struct {
	ID        uint      `json:"id"`
	ToID      uint      `json:"to_id"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Username  string    `json:"username"`
	Photo     string    `json:"photo"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}
type GetRecordListItemList []GetRecordListItem

// GetRecordListResponse 返回用户浏览记录
func GetRecordListResponse(rl *record.RecordList) (data interface{}, err error) {
	list := make(GetRecordListItemList, 0)
	for _, v := range *rl {
		var cover string
		var photo string
		var title string
		var username string
		var tp string
		var toId uint
		if v.Type == "video" {
			cover, _ = conversion.FormattingJsonSrc(v.VideoInfo.Cover)
			photo, _ = conversion.FormattingJsonSrc(v.VideoInfo.UserInfo.Photo)
			title = v.VideoInfo.Title
			username = v.VideoInfo.UserInfo.Username
			tp = "视频"
			toId = *v.ToVideoId
		} else if v.Type == "article" {
			cover, _ = conversion.FormattingJsonSrc(v.ArticleInfo.Cover)
			photo, _ = conversion.FormattingJsonSrc(v.ArticleInfo.UserInfo.Photo)
			title = v.ArticleInfo.Title
			username = v.ArticleInfo.UserInfo.Username
			tp = "专栏"
			toId = *v.ToArticleId
		} else {
			cover, _ = conversion.FormattingJsonSrc(v.UserInfo.LiveInfo.Img)
			photo, _ = conversion.FormattingJsonSrc(v.UserInfo.Photo)
			title = v.UserInfo.LiveInfo.Title
			username = v.UserInfo.Username
			tp = "直播"
			toId = *v.ToLiveId
		}
		list = append(list, GetRecordListItem{
			ID:        v.ID,
			ToID:      toId,
			Title:     title,
			Cover:     cover,
			Username:  username,
			Photo:     photo,
			Type:      tp,
			UpdatedAt: v.UpdatedAt,
		})
	}
	return list, nil
}

// GetNoticeListItem 获取通知列表
type GetNoticeListItem struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	CID       *uint     `json:"cid"`
	Type      string    `json:"type"`
	ToID      *uint     `json:"to_id"`
	Photo     string    `json:"photo"`
	Comment   string    `json:"comment"`
	Cover     string    `json:"cover"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type GetNoticeListStruct []GetNoticeListItem

// GetNoticeListResponse 获取通知列表
func GetNoticeListResponse(nl *notice.NoticesList) (data interface{}, err error) {
	list := make(GetNoticeListStruct, 0)
	for _, v := range *nl {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo)
		var cover string
		var title string

		//判断类型确定标题和封面
		switch v.Type {
		case notice.VideoComment:
			cover, _ = conversion.FormattingJsonSrc(v.VideoInfo.Cover)
			title = v.VideoInfo.Title
			break
		case notice.VideoLike:
			cover, _ = conversion.FormattingJsonSrc(v.VideoInfo.Cover)
			title = v.VideoInfo.Title
			break
		case notice.ArticleComment:
			cover, _ = conversion.FormattingJsonSrc(v.ArticleInfo.Cover)
			title = v.ArticleInfo.Title
			break
		case notice.ArticleLike:
			cover, _ = conversion.FormattingJsonSrc(v.ArticleInfo.Cover)
			title = v.ArticleInfo.Title
			break
		}

		list = append(list, GetNoticeListItem{
			ID:        v.ID,
			Type:      v.Type,
			ToID:      v.ToID,
			CID:       v.Cid,
			Username:  v.UserInfo.Username,
			Photo:     photo,
			Comment:   v.Content,
			Cover:     cover,
			Title:     title,
			CreatedAt: v.CreatedAt,
		})

	}

	return list, nil
}

// ChatMessageInfo 聊天信息
type ChatMessageInfo struct {
	ID        uint      `json:"id"`
	Uid       uint      `json:"uid"`
	Username  string    `json:"username"`
	Photo     string    `json:"photo"`
	Tid       uint      `json:"tid"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// GetChatListItem 聊天列表
type GetChatListItem struct {
	ToID            uint              `json:"to_id"`
	Username        string            `json:"username"`
	Photo           string            `json:"photo"`
	Unread          int               `json:"unread" gorm:"unread"`
	LastMessage     string            `json:"last_message"`
	LastMessagePage int               `json:"last_message_page"`
	MessageList     []ChatMessageInfo `json:"message_list"`
	LastAt          time.Time         `json:"last_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type GetChatListResponseStruct []GetChatListItem

func GetChatListResponse(chatList *chat.ChatList, msgList map[uint]*chat.MsgList) (data interface{}, err error) {
	list := make(GetChatListResponseStruct, 0)
	for _, v := range *chatList {
		photo, _ := conversion.FormattingJsonSrc(v.ToUserInfo.Photo)
		messageList := make([]ChatMessageInfo, 0)
		for _, vv := range *msgList[v.Tid] {
			uPhoto, _ := conversion.FormattingJsonSrc(vv.UInfo.Photo)
			messageList = append(messageList, ChatMessageInfo{
				ID:        vv.ID,
				Uid:       vv.Uid,
				Username:  vv.UInfo.Username,
				Photo:     uPhoto,
				Tid:       vv.Tid,
				Message:   vv.Message,
				Type:      vv.Type,
				CreatedAt: vv.CreatedAt,
			})
		}
		list = append(list, GetChatListItem{
			ToID:        v.Tid,
			Username:    v.ToUserInfo.Username,
			Photo:       photo,
			Unread:      v.Unread,
			LastMessage: v.LastMessage,
			MessageList: messageList,
			LastAt:      v.LastAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}
	return list, nil
}

func GetChatHistoryMsgResponse(list *chat.MsgList) (data interface{}, err error) {
	messageList := make([]ChatMessageInfo, 0)
	for _, v := range *list {
		photo, _ := conversion.FormattingJsonSrc(v.UInfo.Photo)
		messageList = append(messageList, ChatMessageInfo{
			ID:        v.ID,
			Uid:       v.Uid,
			Username:  v.UInfo.Username,
			Photo:     photo,
			Tid:       v.Tid,
			Message:   v.Message,
			Type:      v.Type,
			CreatedAt: v.CreatedAt,
		})
	}
	return messageList, nil
}
