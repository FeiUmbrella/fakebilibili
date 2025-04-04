package live

import (
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/conversion"
	"fakebilibili/infrastructure/pkg/utils/response"
	"fakebilibili/infrastructure/proto/pb"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"strconv"
)

// UserChannel 用户信息
type UserChannel struct {
	UserInfo user.User       // 用户个人信息
	Socket   *websocket.Conn // 前端用户通信的ws
	MsgList  chan []byte     // 给用户发送的信息
}

// LiveRoomEvent 直播间事件
type LiveRoomEvent struct {
	RoomID  uint         // 上线用户要进的房间号
	Channel *UserChannel // 上线用户信息
}

type UserMapChannel map[uint]*UserChannel

type Engin struct {
	LiveRoom map[uint]UserMapChannel // 直播间对应的若干观众

	Register     chan LiveRoomEvent // 上线注册
	Cancellation chan LiveRoomEvent // 下线撤销
}

// Severe 全局变量
var Severe = &Engin{
	LiveRoom:     make(map[uint]UserMapChannel, 10),
	Register:     make(chan LiveRoomEvent, 10),
	Cancellation: make(chan LiveRoomEvent, 10),
}

// Start 启动服务
func (e *Engin) Start() {
	// 为每个用户创建直播间
	type userList []user.User
	users := new(userList)
	global.MysqlDb.Select("id").Find(&users)
	for _, uid := range *users {
		e.LiveRoom[uid.ID] = make(UserMapChannel, 10)
	}

	for {
		select {
		// 注册事件
		case register := <-e.Register:
			// 不能存在该直播间，前端传入直播间房号错误
			if _, ok := e.LiveRoom[register.RoomID]; !ok {
				message := &pb.Message{
					MsgType: consts.Error,
					Data:    []byte("消息格式错误"),
				}
				res, _ := proto.Marshal(message)
				_ = register.Channel.Socket.WriteMessage(websocket.BinaryMessage, res)
				return
			}
			// 添加register到对应直播间
			e.LiveRoom[register.RoomID][register.Channel.UserInfo.ID] = register.Channel
			// 广播用户上线
			err := serviceOnlineAndOfflineRemind(register, true)
			if err != nil {
				response.ErrorWs(register.Channel.Socket, err.Error())
			}
			// 给用户发送直播间历史信息/弹幕
			err = serviceResponseLiveRoomHistoricalBarrage(register)
			if err != nil {
				response.ErrorWs(register.Channel.Socket, err.Error())
			}
		case cancellation := <-e.Cancellation:
			delete(e.LiveRoom[cancellation.RoomID], cancellation.Channel.UserInfo.ID)
			// 广播用户下线
			err := serviceOnlineAndOfflineRemind(cancellation, false)
			if err != nil {
				response.ErrorWs(cancellation.Channel.Socket, err.Error())
			}
		}
	}
}

func CreateSocket(ctx *gin.Context, uid uint, roomID uint, conn *websocket.Conn) error {
	//创建UserChannel
	userChannel := new(UserChannel)
	//绑定ws
	userChannel.Socket = conn
	//绑定用户信息
	u := user.User{}
	u.Find(uid)
	userChannel.UserInfo = u
	//防止阻塞
	userChannel.MsgList = make(chan []byte, 10)

	//创建用户
	userLiveEvent := LiveRoomEvent{
		RoomID:  roomID,
		Channel: userChannel,
	}
	Severe.Register <- userLiveEvent

	go userLiveEvent.Read()
	go userLiveEvent.Write()
	return nil
}

// Write 监听写入数据向前端推送
func (lre LiveRoomEvent) Write() {
	for {
		select {
		case msg := <-lre.Channel.MsgList:
			err := lre.Channel.Socket.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				return
			}
		}
	}
}

// Read 监听前端发送来的弹幕
func (lre LiveRoomEvent) Read() {
	//链接断开进行离线
	defer func() {
		Severe.Cancellation <- lre
		err := lre.Channel.Socket.Close()
		if err != nil {
			return
		}
	}()

	// 读取前端发送的弹幕
	for {
		lre.Channel.Socket.PongHandler()

		_, message, err := lre.Channel.Socket.ReadMessage()
		if err != nil {
			return
		}
		data := &pb.Message{}
		if err := proto.Unmarshal(message, data); err != nil {
			response.ErrorWsProto(lre.Channel.Socket, "消息格式错误")
		}
		// 处理弹幕
		err = getTypeCorrespondingFunc(lre, data)
		if err != nil {
			response.ErrorWsProto(lre.Channel.Socket, err.Error())
		}
	}
}

func getTypeCorrespondingFunc(lre LiveRoomEvent, data *pb.Message) error {
	switch data.MsgType {
	case consts.WebClientBarrageReq:
		// 前端用户发送的弹幕
		return serviceSendBarrage(lre, data.Data)
	}
	response.ErrorWsProto(lre.Channel.Socket, "未定义的消息格式")
	return nil
}

// serviceSendBarrage 发送弹幕信息
func serviceSendBarrage(lre LiveRoomEvent, text []byte) error {
	barrageInfo := &pb.WebClientSendBarrageReq{}
	if err := proto.Unmarshal(text, barrageInfo); err != nil {
		return fmt.Errorf("消息格式错误")
	}
	src, _ := conversion.FormattingJsonSrc(lre.Channel.UserInfo.Photo)
	resp := &pb.WebClientSendBarrageRes{
		UserId:   float32(lre.Channel.UserInfo.ID),
		Username: lre.Channel.UserInfo.Username,
		Avatar:   src,
		Text:     barrageInfo.Text,
		Color:    barrageInfo.Color,
		Type:     barrageInfo.Type,
	}
	data, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}

	// 将弹幕存入Redis中的最近信息
	str := string(data)
	// 如果近期redis保存该直播间弹幕数量大于20，从列表右边弹出一条弹幕
	if studentLen, _ := global.RedisDb.LLen(consts.LiveRoomHistoricalBarrage + strconv.Itoa(int(lre.RoomID))).Result(); studentLen >= 20 {
		err := global.RedisDb.RPop(consts.LiveRoomHistoricalBarrage + strconv.Itoa(int(lre.RoomID))).Err()
		if err != nil {
			return err
		}
	}
	// 从列表左边插入该条弹幕
	err = global.RedisDb.LPush(consts.LiveRoomHistoricalBarrage+strconv.Itoa(int(lre.RoomID)), str).Err()
	if err != nil {
		global.Logger.Errorf("房间ID： %d 最近弹幕存入Redis失败 消息： %s", lre.RoomID, data)
		return err
	}

	// 格式化响应
	message := &pb.Message{
		MsgType: consts.WebClientBarrageRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	// 遍历直播间所有观众，将该条弹幕推送给观众
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}
	return nil
}

// serviceOnlineAndOfflineRemind 广播用户上/下线
func serviceOnlineAndOfflineRemind(lre LiveRoomEvent, flag bool) error {
	//得到当前所有用户
	type userListStruct []*pb.EnterLiveRoom
	userList := make(userListStruct, 0)
	src, _ := conversion.FormattingJsonSrc(lre.Channel.UserInfo.Photo)
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		itemSrc, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo)
		item := &pb.EnterLiveRoom{
			UserId:   float32(v.UserInfo.ID),
			Username: v.UserInfo.Username,
			Avatar:   itemSrc,
		}
		userList = append(userList, item)
	}
	resp := &pb.WebClientEnterLiveRoomRes{
		UserId:   float32(lre.Channel.UserInfo.ID),
		Username: lre.Channel.UserInfo.Username,
		Avatar:   src,
		Type:     flag,
		List:     userList,
	}

	//响应输出
	data, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	message := &pb.Message{
		MsgType: consts.WebClientEnterLiveRoomRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}
	return nil
}

// serviceResponseLiveRoomHistoricalBarrage 给新进直播间的用户推送历史信息/弹幕
func serviceResponseLiveRoomHistoricalBarrage(lre LiveRoomEvent) error {
	//得到历史消息
	val, err := global.RedisDb.LRange(consts.LiveRoomHistoricalBarrage+strconv.Itoa(int(lre.RoomID)), 0, -1).Result()

	if err != nil {
		return fmt.Errorf("获取历史弹幕失败")
	}
	historicalBarrage := &pb.WebClientHistoricalBarrageRes{}
	list := make([]*pb.WebClientSendBarrageRes, 0)
	for _, v := range val {
		barrage := &pb.WebClientSendBarrageRes{}
		// todo:这里没有使用unsafe.pointer去进行0拷贝转化v
		if err := proto.Unmarshal([]byte(v), barrage); err != nil {
			return fmt.Errorf("消息格式错误")
		}
		list = append(list, barrage)
	}
	historicalBarrage.List = list
	data, err := proto.Marshal(historicalBarrage)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	//格式化响应
	message := &pb.Message{
		MsgType: consts.WebClientHistoricalBarrageRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}

	return nil
}
