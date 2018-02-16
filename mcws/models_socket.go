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

func (uWs *WSocket) Close() error {
	return uWs.Ws.Close()
}

func (uWs *WSocket) SendMsg(msg []byte) error {
	return uWs.Ws.WriteMessage(websocket.TextMessage, msg)
}
