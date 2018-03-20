package mcws

import (
	"github.com/gorilla/websocket"
	"errors"
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

type WSocketsManagerSt struct {
	WSocketsListById map[string]*WSocket
	WSocketRoomsList map[string]*WSocketRoom
}

func NewWSocketsManager() *WSocketsManagerSt {
	return &WSocketsManagerSt{
		map[string]*WSocket{},
		map[string]*WSocketRoom{},
	}
}

func (this *WSocketsManagerSt) AddWSocketById(ws *WSocket) error {
	this.WSocketsListById[ws.Id] = ws
	return nil
}

func (this *WSocketsManagerSt) RemoveWSocketById(id string) error {
	delete(this.WSocketsListById, id)
	return nil
}

func (this *WSocketsManagerSt) SendWSocketMsgById(id string, msg []byte) error {
	return this.WSocketsListById[id].SendMsg(msg)
}

// MANAGE ROOMS

func (this *WSocketsManagerSt) GetOrAddWSocketRoomWithWSocket(roomName string, ws *WSocket) *WSocketRoom {
	if r, ok := this.WSocketRoomsList[roomName]; ok {
		r.GetOrAddWSocket(ws)
		return r
	}

	r := CreateWSocketRoom(roomName)
	r.GetOrAddWSocket(ws)
	this.addWSocketRoom(roomName, r)

	return r
}

func (this *WSocketsManagerSt) addWSocketRoom(roomName string, wsr *WSocketRoom) {
	this.WSocketRoomsList[roomName] = wsr
}

func (this *WSocketsManagerSt) RemoveWSocketRoom(roomName string) {
	delete(this.WSocketRoomsList, roomName)
}

func (this *WSocketsManagerSt) AddWSocketToRoom(roomName string, id string, ws *WSocket) {
	this.WSocketRoomsList[roomName].addWSocket(ws)
}

func (this *WSocketsManagerSt) RemoveWSocketFromRoom(roomName string, id string) {
	this.WSocketRoomsList[roomName].RemoveWSocketById(id)
}

func (this *WSocketsManagerSt) SendWsMsgByRoomName(roomName string, msg WSocketMsgBaseI) error {
	room := this.WSocketRoomsList[roomName]
	if room != nil {
		return this.WSocketRoomsList[roomName].SendMsg(msg)
	} else {
		return errors.New("room not found")
	}
}