package mcws

import (
	"errors"
)

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
	// ToDo: Add mutex
	if _, ok := this.WSocketsListById[id]; ok {
		delete(this.WSocketsListById, id)
	}
	if _, ok := this.WSocketRoomsList[id]; ok {
		delete(this.WSocketsListById, id)
	}
	for key, v := range this.WSocketRoomsList {
		if _, ok := v.WSocketsById[id]; ok {
			delete(v.WSocketsById, id)
			if len(v.WSocketsById) == 0 {
				delete(this.WSocketRoomsList, key)
			}
		}
	}
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