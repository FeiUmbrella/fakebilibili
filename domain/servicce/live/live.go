package live

import (
	"encoding/json"
	"fakebilibili/adapter/http/receive/live"
	live2 "fakebilibili/adapter/http/response/live"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/record"
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// GetLiveRoom 返回开播时对应直播间推流地址和推流码
func GetLiveRoom(uid uint) (interface{}, error) {
	// 请求直播服务器 http://127.0.0.1:8090/control/get?room=${uid}
	// todo:注意这里的推流地址的格式，后面配置媒体服务器可能要用
	url := global.Config.LiveConfig.Agreement + "://" + global.Config.LiveConfig.IP + ":" +
		global.Config.LiveConfig.Api + "/control/get?room=" + strconv.Itoa(int(uid))
	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// 将回复中body的数据解析到定义结构体中
	ReqGetRoom := new(live.ReqGetRoom)
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(ReqGetRoom); err != nil {
		return nil, err
	}
	if ReqGetRoom.Status != 200 {
		return nil, fmt.Errorf("获取直播地址失败")
	}
	// 推流地址rtmp://127.0.0.1:1935/live  推流码：ReqGetRoom.Data
	return live2.GetLiveRoomResponse("rtmp://"+global.Config.LiveConfig.IP+
		":"+global.Config.LiveConfig.RTMP+"/live", ReqGetRoom.Data), nil
}

// GetLiveRoomInfo 给前端返回直播间信息及直播间拉流地址
func GetLiveRoomInfo(data live.GetLiveRoomInfoReceiveStruct, uid uint) (interface{}, error) {
	userInfo := new(user.User)
	userInfo.FindLiveInfo(data.RoomID) // 查询直播间信息，直播间ID就是主播的uid

	// todo:注意这里拉流的地址，后面配置媒体服务器可能要用
	//拉流地址 http://8.138.149.242:7001/live/37.flv
	flv := global.Config.LiveConfig.Agreement + "://" + global.Config.LiveConfig.IP + ":" +
		global.Config.LiveConfig.FLV + "/live/" + strconv.Itoa(int(data.RoomID)) + ".flv"

	if uid > 0 {
		// 添加观看直播的历史记录
		rd := new(record.Record)
		err := rd.AddLiveRecord(uid, data.RoomID)
		if err != nil {
			return nil, fmt.Errorf("添加观看直播历史记录失败")
		}
	}
	return live2.GetLiveRoomInfoResponse(userInfo, flv), nil
}

func GetBeLiveList() (interface{}, error) {
	// http://127.0.0.1:8090//stat/livestat
	// todo：这里从媒体服务器中获取所有publisher的信息的访问地址，后面配置媒体服务器可能要用
	url := global.Config.LiveConfig.Agreement + "://" +
		global.Config.LiveConfig.IP + ":" +
		global.Config.LiveConfig.Api + "/stat/livestat"

	// 创建http get请求，从媒体服务器中获取推流的publisher的信息
	resp, err := http.Get(url)
	global.Logger.Infof("获取直播在线list的请求url为%s,请求返回结果为%v", url, resp)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// 将回复中body的数据解析到定义结构体中
	livestat := new(live.LivestatRes)
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(livestat); err != nil {
		return nil, err
	}
	if livestat.Status != 200 {
		return nil, fmt.Errorf("获取直播列表失败")
	}

	keys := make([]uint, 0)
	for _, kv := range livestat.Data.Publishers {
		// todo:这里获取到的kv是以 live/${id} 结尾的，这里的id就是publisher的uid，配置媒体服务器可能要用
		k := strings.Split(kv.Key, "live/")
		uintKey, _ := strconv.ParseUint(k[1], 10, 19)
		keys = append(keys, uint(uintKey))
	}
	global.Logger.Infof("查询userList的keys为%v", keys)
	userList := new(user.UserList)
	if len(keys) > 0 {
		err = userList.GetBeLiveList(keys)
		if err != nil {
			return nil, fmt.Errorf("查询失败")
		}
	}
	return live2.GetBeLiveListResponse(userList), nil
}
