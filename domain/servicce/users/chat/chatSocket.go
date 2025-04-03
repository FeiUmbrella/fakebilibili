package chat

import (
	"encoding/json"
	"fakebilibili/adapter/http/receive/socket"
	"fakebilibili/domain/servicce/users"
	"fakebilibili/infrastructure/consts"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/chat"
	"fakebilibili/infrastructure/pkg/utils/response"
	"fmt"
	"github.com/gorilla/websocket"
)

// MsgInfo 私信的类型和内容
type MsgInfo struct {
	Type string
	Data interface{}
}

// UserChannel 包含a用户信息、与a用户通信的ws、以及存放给a用户发送私信的chan
type UserChannel struct {
	UserInfo *user.User
	Socket   *websocket.Conn
	ChatList map[uint]*websocket.Conn // todo:存储当前在a用户“私信界面列表”的用户的ws?
	MegList  chan MsgInfo             // 给用户推送的信息
}

type Engine struct {
	UserMapChannel map[uint]*UserChannel // 保存所有当前在“私信界面”的用户的UserChannel. uid -> UserChannel
	Register       chan *UserChannel     // 上线
	Cancellation   chan *UserChannel     // 下线
}

// Severe 全局变量
var Severe = &Engine{
	UserMapChannel: make(map[uint]*UserChannel, 10),
	Register:       make(chan *UserChannel, 10),
	Cancellation:   make(chan *UserChannel, 10),
}

// Start 启动推送通知的服务
func (e *Engine) Start() {
	for {
		select {
		case register := <-e.Register: // 有新成员上线
			e.UserMapChannel[register.UserInfo.ID] = register // 添加新成员
			// 向该上线用户推送未读私信
			register.ChatMessage(consts.ChatOnlineUnreadMsg)
		case cancellation := <-e.Cancellation: // 有成员断开连接下线
			//fmt.Printf("%d成员下线", cancellation.UserInfo.ID)
			delete(e.UserMapChannel, cancellation.UserInfo.ID) // 删除该成员
		}
	}
}

// CreateChatSocket 创建监听信息和发送信息的chat-Ws
func CreateChatSocket(uid uint, conn *websocket.Conn) error {
	// 创建一个新的包含用户信息的UserChannel
	userChannel := new(UserChannel)
	// 将传入的ws绑定到该用户的userChannel
	userChannel.Socket = conn
	// 将该用户的信息绑定到该用户的userChannel
	u := new(user.User)
	u.Find(uid)
	userChannel.UserInfo = u
	// 初始化
	userChannel.MegList = make(chan MsgInfo, 10)
	userChannel.ChatList = make(map[uint]*websocket.Conn)
	// 该用户上线
	Severe.Register <- userChannel

	// 并发利用ws监听用户发送的信息已经利用ws给用户发信息
	go userChannel.Read()
	go userChannel.Write()
	return nil
}

// Write 利用ws向前端用户推送信息
func (uc *UserChannel) Write() {
	for { // 当前端断开ws连接的时候，自动跳出for循环
		select { // 一直监听channel，有要发送给用户的信息就取出推送给前端用户
		case msg := <-uc.MegList:
			//fmt.Println("推送！")
			response.SuccessWs(uc.Socket, msg.Type, msg.Data)
		}
	}
}

// Read 利用ws监听前端用户发送的信息
func (uc *UserChannel) Read() {
	// 用户退出网站，离线断开ws连接
	defer func() {
		Severe.Cancellation <- uc // 将该用户放入离线channel
		err := uc.Socket.Close()
		if err != nil {
			return
		}
	}()
	for { // 当前端断开ws连接的时候，自动跳出for循环，接着运行上面的延迟函数
		// ping客户端检测客户端是否在线，如果在线客户端会返回pong，服务端接收到pong说明客户端在线
		uc.Socket.PongHandler()                 // 接收到Pong时使用默认处理函数-不进行任何处理
		_, text, err := uc.Socket.ReadMessage() // 读取发送给后端的信息
		if err != nil {
			return
		}
		//fmt.Println("接收到前端数据")
		info := new(socket.Receive)
		if err = json.Unmarshal(text, info); err != nil {
			response.ErrorWs(uc.Socket, "发送的消息格式错误")
		}
		switch info.Type {
		// todo:这里接收到信息后没有做任何处理和相应，说明重点是后端向前端推送信息，前端不会向后端发送信息或者说前端发送的信息没用
		}
	}
}

// ChatMessage 推送离线时未读私信
func (uc *UserChannel) ChatMessage(msgType string) {
	cl := new(chat.ChatsListInfo)
	unreadNum := cl.GetUnreadNumber(uc.UserInfo.ID)
	if *unreadNum > 0 {
		// 存在离线时接收的未读私信,直接推送已经存在的私信列表与聊天信息
		//fmt.Printf("未读私信数量为：%d\n", *unreadNum)
		list, err := users.GetChatList(uc.UserInfo.ID)
		if err != nil {
			fmt.Println("查询聊天记录错误")
			return
		}
		// 推送
		uc.MegList <- MsgInfo{msgType, list}
	}
}
