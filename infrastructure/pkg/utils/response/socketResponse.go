package response

import (
	"fakebilibili/infrastructure/consts"
	"github.com/gorilla/websocket"
)

type DataWs struct {
	Code    MyCode      `json:"code"`
	Type    string      `json:"type,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // 如果没有传入该参数则返回时省略该参数
}

// NotLoginWs 没有成功解析用户身份
func NotLoginWs(ws *websocket.Conn, msg string) {
	resData := &DataWs{
		Code:    CodeNotLogin,
		Message: msg,
		Data:    nil,
	}
	err := ws.WriteJSON(resData) // 通过ws主动向前端发送信息
	if err != nil {
		return
	}
}

// SuccessWs 利用ws向前端发送成功信息
func SuccessWs(ws *websocket.Conn, tp string, data interface{}) {
	resData := &DataWs{
		Code:    CodeSuccess,
		Type:    tp,
		Message: CodeSuccess.Msg(),
		Data:    data,
	}
	err := ws.WriteJSON(resData)
	if err != nil {
		return
	}
}

// ErrorWs 利用ws向前端发送错误信息
func ErrorWs(ws *websocket.Conn, msg string) {
	resData := &DataWs{
		Code:    CodeServerBusy,
		Type:    consts.VideoSocketTypeError,
		Message: msg,
		Data:    nil,
	}
	err := ws.WriteJSON(resData)
	if err != nil {
		return
	}
}
