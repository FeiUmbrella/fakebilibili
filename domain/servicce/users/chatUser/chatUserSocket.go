package chatUser

import (
	"encoding/json"
	"fakebilibili/adapter/http/receive/socket"
	chat2 "fakebilibili/domain/servicce/users/chat"
	"fakebilibili/infrastructure/model/user"
	"fakebilibili/infrastructure/model/user/chat"
	"fakebilibili/infrastructure/pkg/global"
	"fakebilibili/infrastructure/pkg/utils/response"
	"github.com/FeiUmbrella/Sensitive_Words_Filter/filter"
	"github.com/gorilla/websocket"
)

// MsgInfo 聊天的类型和内容
type MsgInfo struct {
	Type string
	Data interface{}
}

// UserChannel 包含a用户信息、与a用户通信的ws、以及存放给a用户发送私信的chan
type UserChannel struct {
	UserInfo *user.User
	Tid      uint // a用户聊天的对象b用户
	Socket   *websocket.Conn
	MegList  chan MsgInfo // 给用户推送的信息
}

type Engine struct {
	UserMapChannel map[uint]*UserChannel // 保存所有当前在聊天页面的用户的UserChannel. uid -> UserChannel
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
			// uid用户此时在tid用户的聊天界面，清空uid用户与tid用户之间未读信息
			chatInfo := new(chat.ChatsListInfo)
			err := chatInfo.UnreadEmpty(register.UserInfo.ID, register.Tid)
			if err != nil {
				global.Logger.Errorf("uid:%d 清空聊天列表中对应 tid:%d 的未读状态失败", register.UserInfo.ID, register.Tid)
			}

			// 此时uid的ChatListWs肯定是在线的
			if _, ok := chat2.Severe.UserMapChannel[register.UserInfo.ID]; ok {
				// 将这里的uid与tid聊天界面的uid的ws注册到ChatList那里，表明uid用户正处在tid用户聊天界面
				chat2.Severe.UserMapChannel[register.UserInfo.ID].ChatList[register.Tid] = register.Socket
			}
		case cancellation := <-e.Cancellation: // 有成员断开连接下线
			delete(e.UserMapChannel, cancellation.UserInfo.ID) // 删除该成员
			if _, ok := chat2.Severe.UserMapChannel[cancellation.UserInfo.ID]; ok {
				// 将注册在ChatList那里的uid的ChatUserSocketWs删除，uid用户退出tid用户聊天界面
				delete(chat2.Severe.UserMapChannel[cancellation.UserInfo.ID].ChatList, cancellation.Tid)
			}
		}
	}
}
func CreateChatByUserSocket(uid, tid uint, conn *websocket.Conn) error {
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
	userChannel.Tid = tid
	// 该用户上线
	Severe.Register <- userChannel

	// 并发利用ws监听用户发送的信息以及利用ws给用户发信息
	go userChannel.Read()
	go userChannel.Write()
	return nil
}

// Writer 监听写入数据
func (uc *UserChannel) Write() {
	for {
		select {
		case msg := <-uc.MegList:
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
		node := filter.NewNode()
		node.StartMatch(info.Data)
		if node.IsSensitive() {
			response.ErrorWs(uc.Socket, "消息中含有敏感词汇！")
			return
		}
		switch info.Type {
		case "sendChatMsgText":
			// 给tid发送信息
			sendChatMsgText(uc, uc.UserInfo.ID, uc.Tid, info)
		}
	}
}
