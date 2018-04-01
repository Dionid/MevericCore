package mcws

import (
	"errors"
	"github.com/gorilla/websocket"
)

type WSocketRoom struct {
	Name string
	WSocketsById map[string]*WSocket
}

func CreateWSocketRoom (name string) *WSocketRoom {
	return &WSocketRoom{
		Name:       name,
		WSocketsById: map[string]*WSocket{},
	}
}

func (r *WSocketRoom) GetOrAddWSocket(ws *WSocket) *WSocket {
	if w, ok := r.WSocketsById[ws.Id]; ok {
		w.Ws = ws.Ws
		return w
	}
	r.WSocketsById[ws.Id] = ws
	return ws
}

func (r *WSocketRoom) addWSocket(ws *WSocket) {
	r.WSocketsById[ws.Id] = ws
}

func (r *WSocketRoom) RemoveWSocketById(id string) {
	delete(r.WSocketsById, id)
}

func (r *WSocketRoom) Close() error {
	for _, ws := range r.WSocketsById {
		if err := ws.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *WSocketRoom) SendMsg(msg WSocketMsgBaseI) (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("recovered: method WsRoom.SendWsMsg()")
		}
	}()

	byteData, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	for _, ws := range r.WSocketsById {
		err = ws.SendMsg(byteData)
		if err != nil {
			if err == websocket.ErrCloseSent {
				continue
			} else {
				return err
			}
		}
	}

	return nil
}
