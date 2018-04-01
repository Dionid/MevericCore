package mcws

import "github.com/gorilla/websocket"

type WSocket struct {
	Id string
	Authorized bool
	Ws *websocket.Conn
}

func CreateWSocket(id string, ws *websocket.Conn) *WSocket {
	return &WSocket{
		id,
		false,
		ws,
	}
}

func (ws *WSocket) Close() error {
	return ws.Ws.Close()
}

//func (ws *WSocket) IsClosed() error {
//	return ws.Ws.
//}

//func (ws *WSocket) SendErrorMsg(err string, action string, errorCode int, reqId int) error {
//	errMsg := CreateWsResActionSingleErrorMsg("Token is required", action, 503, reqId)
//	msg, _ := errMsg.MarshalJSON()
//	return ws.SendMsg(msg)
//}

func (ws *WSocket) SendMsg(msg []byte) error {
	return ws.Ws.WriteMessage(websocket.TextMessage, msg)
}
