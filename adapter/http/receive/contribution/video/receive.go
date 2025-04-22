package video

import "fakebilibili/infrastructure/model/common"

type GetVideoBarrageReceiveStruct struct {
	ID string `json:"id"`
}

type GetVideoBarrageListReceiveStruct struct {
	ID string `json:"id"`
}

type GetVideoCommentReceiveStruct struct {
	PageInfo common.PageInfo `json:"pageInfo"`
	VideoID  uint            `json:"video_id" binding:"required"`
}

type LikeVideoCommentReqStruct struct {
	VideoCommentId int  `json:"video_comment_id"`
	VideoId        uint `json:"video_id"`
}

type GetVideoContributionByIDReceiveStruct struct {
	VideoID uint `json:"video_id"`
}

type SendVideoBarrageReceiveStruct struct {
	Author string  `json:"author"`
	Color  uint    `json:"color" binding:"required"`
	ID     string  `json:"id" binding:"required"`
	Text   string  `json:"text" binding:"required"`
	Time   float64 `json:"time"`
	Type   uint    `json:"type"`
	Token  string  `json:"token" binding:"required"`
}

type CreateVideoContributionReceiveStruct struct {
	VideoPath       string   `json:"video" binding:"required"` // 视频文件路径
	VideoUploadType string   `json:"videoUploadType" binding:"required"`
	CoverPath       string   `json:"cover" binding:"required"`
	CoverUploadType string   `json:"coverUploadType" binding:"required"`
	Title           string   `json:"title" binding:"required"`
	Reprinted       *bool    `json:"reprinted" binding:"required"`
	Label           []string `json:"label"`
	Introduce       string   `json:"introduce" binding:"required"`
	VideoDuration   int64    `json:"videoDuration" binding:"required"`
	Media           *string  `json:"media"` // 将oss上视频注册媒体资源后的到的对应媒体资源ID
	DateTime        string   `json:"date1time"`
}

type UpdateVideoContributionReceiveStruct struct {
	ID              uint     `json:"id" binding:"required"`
	Cover           string   `json:"cover" binding:"required"`
	CoverUploadType string   `json:"coverUploadType" binding:"required"`
	Title           string   `json:"title" binding:"required"`
	Reprinted       *bool    `json:"reprinted" binding:"required"`
	Label           []string `json:"label"`
	Introduce       string   `json:"introduce" binding:"required"`
}

type DeleteVideoByIDReceiveStruct struct {
	ID uint `json:"id"`
}

type VideosPostCommentReceiveStruct struct {
	VideoID   uint   `json:"video_id"`
	Content   string `json:"content"`
	ContentID uint   `json:"content_id"`
}

type GetVideoManagementListReceiveStruct struct {
	PageInfo common.PageInfo `json:"page_info"`
}

type LikeVideoReceiveStruct struct {
	VideoID uint `json:"video_id"`
}

type DeleteVideoByPathReceiveStruct struct {
	Path string `json:"path"`
}

type SendWatchTimeReqStruct struct {
	Id   string `json:"id"`
	Time string `json:"time"`
}
