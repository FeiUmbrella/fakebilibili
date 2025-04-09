package live

import (
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/utils/conversion"
)

type GetLiveRoomResponseStruct struct {
	Address string `json:"address"` // 推流地址
	Key     string `json:"key"`     // 推流码
}

func GetLiveRoomResponse(address string, key string) interface{} {
	return GetLiveRoomResponseStruct{
		Address: address,
		Key:     key,
	}
}

// GetLiveRoomInfoResponseStruct 返回直播间相关信息包括拉流地址
type GetLiveRoomInfoResponseStruct struct {
	Username  string `json:"username"`
	Photo     string `json:"photo"`
	LiveTitle string `json:"live_title"`
	Flv       string `json:"flv"` // 拉流地址
}

func GetLiveRoomInfoResponse(info *user.User, flv string) interface{} {
	photo, _ := conversion.FormattingJsonSrc(info.Photo)
	return GetLiveRoomInfoResponseStruct{
		Username:  info.Username,
		Photo:     photo,
		LiveTitle: info.LiveInfo.Title,
		Flv:       flv,
	}
}

type BeLiveInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
	Img      string `json:"img"`
	Title    string `json:"title"`
	Online   int    `json:"online"`
}

type BeLiveInfoList []BeLiveInfo

// GetBeLiveListResponse 返回直播列表信息
func GetBeLiveListResponse(ul *user.UserList) interface{} {
	list := make(BeLiveInfoList, 0)
	for _, v := range *ul {
		photo, _ := conversion.FormattingJsonSrc(v.Photo)
		img, _ := conversion.FormattingJsonSrc(v.LiveInfo.Img)
		list = append(list, BeLiveInfo{
			ID:       v.ID,
			Username: v.Username,
			Photo:    photo,
			Img:      img,
			Title:    v.LiveInfo.Title,
			Online:   0,
		})
	}
	return list
}
