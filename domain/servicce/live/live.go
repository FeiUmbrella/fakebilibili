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
	//// 请求直播服务器 http://127.0.0.1:8090/control/get?room=${uid}
	//url := global.Config.LiveConfig.Agreement + "://" + global.Config.LiveConfig.IP + ":" +
	//	global.Config.LiveConfig.Api + "/control/get?room=" + strconv.Itoa(int(uid))
	//// 创建http get请求
	//resp, err := http.Get(url)
	//if err != nil {
	//	return nil, err
	//}
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//
	//	}
	//}(resp.Body)
	//// 将回复中body的数据解析到定义结构体中
	//ReqGetRoom := new(live.ReqGetRoom)
	//decoder := json.NewDecoder(resp.Body)
	//if err := decoder.Decode(ReqGetRoom); err != nil {
	//	return nil, err
	//}
	//if ReqGetRoom.Status != 200 {
	//	return nil, fmt.Errorf("获取直播地址失败")
	//}
	//// 推流地址rtmp://127.0.0.1:1935/live  推流码：ReqGetRoom.Data
	//return live2.GetLiveRoomResponse("rtmp://"+global.Config.LiveConfig.IP+
	//	":"+global.Config.LiveConfig.RTMP+"/live", ReqGetRoom.Data), nil

	// rtmp推流为 rtmp://47.97.31.45:1935/app
	// 推流码为 room-${uid}
	ReqGetRoom := new(live.ReqGetRoom)
	ReqGetRoom.Data = fmt.Sprintf("room-%d", uid)
	address := "rtmp://" + global.Config.LiveConfig.IP + ":" + global.Config.LiveConfig.RTMP + "/app"

	return live2.GetLiveRoomResponse(address, ReqGetRoom.Data), nil
}

// GetLiveRoomInfo 给前端返回直播间信息及直播间拉流地址
func GetLiveRoomInfo(data live.GetLiveRoomInfoReceiveStruct, uid uint) (interface{}, error) {
	userInfo := new(user.User)
	userInfo.FindLiveInfo(data.RoomID) // 查询直播间信息，直播间ID就是主播的uid

	////拉流地址 http://8.138.149.242:7001/live/37.flv
	//flv := global.Config.LiveConfig.Agreement + "://" + global.Config.LiveConfig.IP + ":" +
	//	global.Config.LiveConfig.FLV + "/live/" + strconv.Itoa(int(data.RoomID)) + ".flv"

	// flv拉流地址为 http://47.97.31.45:80/live?port=1935&app=app&stream=room-1
	flv := global.Config.LiveConfig.Agreement + "://" + global.Config.LiveConfig.IP + ":" +
		global.Config.LiveConfig.FLV + "/live?" + "port=" + global.Config.LiveConfig.RTMP +
		"&app=app" + "&stream=room-" +
		strconv.Itoa(int(data.RoomID))

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
	// http://47.97.31.45:80/stat/livestat http.Get得到正在推流的推流码room-${uid}
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取直播列表失败")
	}
	// 将回复中body的数据解析到定义结构体中
	livestat := new(live.LivestatRes)
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(livestat); err != nil {
		return nil, err
	}

	keys := make([]uint, 0)
	for _, server := range livestat.HTTPFLV.Servers {
		for _, app := range server.Applications {
			for _, stream := range app.Live.Streams {
				// stream.Name: room-${uid}
				k := strings.Split(stream.Name, "-")
				uintKey, _ := strconv.ParseUint(k[1], 10, 19)
				keys = append(keys, uint(uintKey))
			}
		}
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
