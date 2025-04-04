package socket

import (
	"encoding/json"
	"fakebilibili/adapter/http/receive/socket"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/pkg/utils/response"
	"github.com/gorilla/websocket"
)

// MsgInfo 信息的类型和内容
type MsgInfo struct {
	Type string
	Data interface{}
}

// UserChannel 包含用户信息、与用户通信的ws、以及存放给该用户发送信息的chan
type UserChannel struct {
	UserInfo user.User
	Socket   *websocket.Conn
	MegList  chan MsgInfo // 给用户推送的信息
}

// UserMapChannel   uid -> *UserChannel对象
type UserMapChannel map[uint]*UserChannel

// VideoRoomEvent 事件 注册 离线
type VideoRoomEvent struct {
	VideoID uint         // 观看视频ID
	Channel *UserChannel // 注册观看视频的用户信息
}
type Engine struct {
	VideoRoom    map[uint]UserMapChannel // 正在观看某一视频的所有用户
	Register     chan VideoRoomEvent     // 上线
	Cancellation chan VideoRoomEvent     // 下线
}

// Severe 全局变量
var Severe = &Engine{
	VideoRoom:    make(map[uint]UserMapChannel, 10),
	Register:     make(chan VideoRoomEvent, 10),
	Cancellation: make(chan VideoRoomEvent, 10),
}

// Start 启动推送通知的服务
func (e *Engine) Start() {
	//fmt.Println("开始Start服务")
	for {
		select {
		case register := <-e.Register: // 有新成员观看该视频
			//添加成员
			e.VideoRoom[register.VideoID][register.Channel.UserInfo.ID] = register.Channel
			// 向观看该视频的所有观众推送当前观看该视频的人数
			register.ViewerNum(consts.VideoSocketTypeNumberOfViewers, e)
		case cancellation := <-e.Cancellation: // 有成员关闭该视频
			// 删除该用户
			delete(e.VideoRoom[cancellation.VideoID], cancellation.Channel.UserInfo.ID)

			// 向观看该视频的所有观众更新当前观看该视频的人数
			cancellation.ViewerNum(consts.VideoSocketTypeNumberOfViewers, e)
		}
	}
}

// CreateVideoSocket 创建监听信息和发送信息的video-Ws
func CreateVideoSocket(uid uint, videoID uint, conn *websocket.Conn) error {
	// 创建一个新的包含用户信息的UserChannel
	userChannel := new(UserChannel)
	// 将传入的ws绑定到该用户的userChannel
	userChannel.Socket = conn
	// 将该用户的信息绑定到该用户的userChannel
	u := new(user.User)
	u.Find(uid)
	userChannel.UserInfo = *u
	// mad，这里忘记初始化MsgList了，导致下面信息放不进去一直阻塞
	userChannel.MegList = make(chan MsgInfo, 20)

	//创建用户
	userLiveEvent := VideoRoomEvent{
		VideoID: videoID,
		Channel: userChannel,
	}
	// 该用户上线
	Severe.Register <- userLiveEvent

	// 并发利用ws监听用户发送的信息已经利用ws给用户发信息
	go userLiveEvent.Read()
	go userLiveEvent.Write()
	return nil
}

// Write 利用ws向用户推送信息
func (vre *VideoRoomEvent) Write() {
	for { // 当前端断开ws连接的时候，自动跳出for循环
		select { // 一直监听channel，有要发送给用户的信息就取出推送给用户
		case msg := <-vre.Channel.MegList:
			//fmt.Println("推送！")
			response.SuccessWs(vre.Channel.Socket, msg.Type, msg.Data)
		}
	}
}

// Read 利用ws监听用户发送的信息
func (vre *VideoRoomEvent) Read() {
	// 用户退出网站，离线断开ws连接
	defer func() {
		Severe.Cancellation <- *vre // 将该用户放入离线channel
		err := vre.Channel.Socket.Close()
		if err != nil {
			return
		}
	}()
	for { // 当前端断开ws连接的时候，自动跳出for循环，接着运行上面的延迟函数
		// ping客户端检测客户端是否在线，如果在线客户端会返回pong，服务端接收到pong说明客户端在线
		vre.Channel.Socket.PongHandler()                 // 接收到Pong时使用默认处理函数-不进行任何处理
		_, text, err := vre.Channel.Socket.ReadMessage() // 读取发送给后端的信息
		if err != nil {
			return
		}
		//fmt.Println("接收到前端数据")
		info := new(socket.Receive)
		if err = json.Unmarshal(text, info); err != nil {
			response.ErrorWs(vre.Channel.Socket, "发送的消息格式错误")
		}
		switch info.Type {
		// todo:这里接收到信息后没有做任何处理和相应，说明重点是后端向前端推送信息，前端不会向后端发送信息或者说前端发送的信息没用
		}
	}
}

func (vre *VideoRoomEvent) ViewerNum(msgType string, e *Engine) {
	num := len(e.VideoRoom[vre.VideoID])
	r := struct {
		People int `json:"people"`
	}{
		People: num,
	}
	res := MsgInfo{
		Type: msgType,
		Data: r,
	}
	// 向观看该视频的所有观众推送当前观看该视频的人数
	for _, v := range e.VideoRoom[vre.VideoID] {
		v.MegList <- res
	}
}
