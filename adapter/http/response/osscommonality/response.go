package osscommonality

import (
	"fakebilibili/adapter/http/receive/osscommonality"
	"fakebilibili/infrastructure/model/contribution/video"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/calculate"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	"time"
)

type GteStsInfoStruct struct {
	Region          string `json:"region"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	StsToken        string `json:"sts_token"`
	Bucket          string `json:"bucket"`
	ExpirationTime  int64  `json:"expiration_time"`
}

// GteStsInfo 返回获取的STS临时授权码信息
func GteStsInfo(info *sts20150401.AssumeRoleResponseBodyCredentials) (interface{}, error) {
	return GteStsInfoStruct{
		Region:          global.Config.AliyunOss.Region,
		AccessKeyID:     *info.AccessKeyId,
		AccessKeySecret: *info.AccessKeySecret,
		StsToken:        *info.SecurityToken,
		Bucket:          global.Config.AliyunOss.Bucket,
		ExpirationTime:  time.Now().Unix() + int64(global.Config.AliyunOss.DurationSeconds),
	}, nil
}

type UploadCheckStruct struct {
	IsUpload bool                           `json:"is_upload"`
	List     osscommonality.UploadSliceList `json:"list"`
	Path     string                         `json:"path"`
}

// UploadCheckResponse 返回文件是否上传，若为上传返回未上传分片列表
func UploadCheckResponse(is bool, list osscommonality.UploadSliceList, path string) (interface{}, error) {
	return UploadCheckStruct{
		IsUpload: is,   // 原文件是否已存在本地
		List:     list, // 未上传在本地的分片列表
		Path:     path, // 若原文件已存在本地，返回本地存储路径
	}, nil
}

func UploadingMethodResponse(tp string) interface{} {
	type UploadingMethodResponseStruct struct {
		Tp string `json:"type"`
	}
	return UploadingMethodResponseStruct{
		Tp: tp,
	}
}

func UploadingDirResponse(dir string, quality float64) interface{} {
	type UploadingDirResponseStruct struct {
		Path    string  `json:"path"`
		Quality float64 `json:"quality"`
	}
	return UploadingDirResponseStruct{
		Path:    dir,
		Quality: quality,
	}
}

// VideoInfo 首页视频
type VideoInfo struct {
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

type videoInfoList []VideoInfo

func SearchVideoResponse(videoList *video.VideosContributionList) (interface{}, error) {
	//处理视频
	vl := make(videoInfoList, 0)
	for _, lk := range *videoList {
		cover, _ := conversion.FormattingJsonSrc(lk.Cover)
		videoSrc, _ := conversion.FormattingJsonSrc(lk.Video)
		info := VideoInfo{
			ID:            lk.ID,
			Uid:           lk.Uid,
			Title:         lk.Title,
			Video:         videoSrc,
			Cover:         cover,
			VideoDuration: lk.VideoDuration,
			Label:         conversion.StringConversionMap(lk.Label),
			Introduce:     lk.Introduce,
			Heat:          lk.Heat,
			BarrageNumber: len(lk.Barrage),
			Username:      lk.UserInfo.Username,
			CreatedAt:     lk.CreatedAt,
		}
		vl = append(vl, info)
	}
	return vl, nil
}

type UserInfo struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Photo       string `json:"photo"`
	Signature   string `json:"signature"`
	IsAttention bool   `json:"is_attention"`
}

type UserInfoList []UserInfo

func SearchUserResponse(userList *user.UserList, aids []uint) (interface{}, error) {
	list := make(UserInfoList, 0)
	for _, v := range *userList {
		photo, _ := conversion.FormattingJsonSrc(v.Photo)
		list = append(list, UserInfo{
			ID:          v.ID,
			Username:    v.Username,
			Photo:       photo,
			Signature:   v.Signature,
			IsAttention: calculate.ArrayIsContain(aids, v.ID), // 查找到的用户v是否是用户的关注
		})
	}
	return list, nil
}
