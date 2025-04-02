package notice

import (
	"fakebilibili/infrastructure/model/common"
	user2 "fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Notice 用户通知（包括：视频/文章评论、点赞、系统通知）
type Notice struct {
	gorm.Model
	Uid     uint   `json:"uid"`                           // 接受notice通知的userId
	Cid     uint   `json:"cid"`                           // 触发这个notice事件的userId
	Type    string `json:"type" gorm:"type:varchar(255)"` // (comment,like,system)
	ToID    uint   `json:"to_id" gorm:"column:to_id"`     // 跳转的视频或文章的id
	ISRead  uint   `json:"is_read" gorm:"column:is_read"`
	Content string `json:"content" gorm:"type:text"`

	VideoInfo   VideoInfo  `json:"videoInfo" gorm:"foreignKey:to_id"`
	UserInfo    user2.User `json:"userInfo" gorm:"foreignKey:cid"`
	ArticleInfo Article    `json:"articleInfo" gorm:"foreignKey:to_id"`
}

var (
	Online         = "online"         //上线时进行通知
	VideoComment   = "videoComment"   //视频评论
	VideoLike      = "videoLike"      //视频点赞
	ArticleComment = "articleComment" //文章评论
	ArticleLike    = "articleLike"    //文章点赞
	UserLogin      = "userLogin"      //用户登录的欢迎消息
	DailyReport    = "dailyReport"    //日报
	UserRegister   = "userRegister"
)

type NoticesList []Notice

func (Notice) TableName() string {
	return "lv_users_notices"
}

type VideoInfo struct {
	gorm.Model
	Uid   uint           `json:"uid"`
	Title string         `json:"title" gorm:"type:varchar(255)"`
	Video datatypes.JSON `json:"video"`
	Cover datatypes.JSON `json:"cover"`
}

func (VideoInfo) TableName() string {
	return "lv_video_contribution"
}

type Article struct {
	gorm.Model
	Uid              uint           `json:"uid"`
	ClassificationID uint           `json:"classification_id"  gorm:"classification_id"`
	Title            string         `json:"title"`
	Cover            datatypes.JSON `json:"cover"`
}

func (Article) TableName() string {
	return "lv_article_contribution"
}

// GetNoticeList 获取通知列表
func (nl *NoticesList) GetNoticeList(page common.PageInfo, msgType []string, uid uint) error {
	if len(msgType) > 0 {
		return global.MysqlDb.
			Where("uid = ?", uid).
			Where("type", msgType).
			Preload("VideoInfo").
			Preload("UserInfo").
			Preload("ArticleInfo").
			Limit(page.Size).
			Offset((page.Page - 1) * page.Size).
			Order("created_at desc").
			Find(nl).Error
	} else {
		return global.MysqlDb.
			Where("uid = ?", uid).
			Preload("VideoInfo").
			Preload("UserInfo").
			Preload("ArticleInfo").
			Limit(page.Size).
			Offset((page.Page - 1) * page.Size).
			Order("created_at desc").
			Find(nl).Error
	}
}

// ReadAll 将所有通知设为已读
func (nt *Notice) ReadAll(uid uint) error {
	return global.MysqlDb.
		Where("uid = ? AND is_read = ?", uid, 0).
		Updates(&Notice{ISRead: 1}).Error
}

// GetUnreadNum 返回未读通知数量
func (nt *Notice) GetUnreadNum(uid uint) *int64 {
	var num int64
	err := global.MysqlDb.Model(&Notice{}).
		Where("uid = ? AND is_read = ?", uid, 0).
		Count(&num).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &num
}

// AddNotice uid：接受notice通知的userId；cid：触发这个notice事件的userId；tid：跳转的视频或文章的id；ty：事件类型；content：通知内容
func (nt *Notice) AddNotice(uid uint, cid uint, tid uint, tp string, content string) error {
	nt.Uid = uid
	nt.Cid = cid
	nt.ToID = tid
	nt.Type = tp
	nt.Content = content
	nt.ISRead = 0
	return global.MysqlDb.Create(nt).Error
}
