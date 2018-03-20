package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
	"strings"
	"gopkg.in/mgo.v2/bson"
)

type DeviceRPCManagerSt struct {
	DeviceResponseServiceSt
	DevicesCollectionManager   mccommon.DevicesCollectionManagerInterface
	DeviceReqControllersByType map[string]DeviceRPCCtrlInterface
}

func CreateDeviceRPCManager(serverId string, devicesColManager mccommon.DevicesCollectionManagerInterface,  SendToUser mccommon.SendRPCMsgFn) *DeviceRPCManagerSt {
	dev := &DeviceRPCManagerSt{
		DevicesCollectionManager: devicesColManager,
		DeviceResponseServiceSt: DeviceResponseServiceSt{
			SendRPCMsgToUser: SendToUser,
			ServerId:         serverId,
		},
		DeviceReqControllersByType: map[string]DeviceRPCCtrlInterface{},
	}
	return dev
}

func (thisR *DeviceRPCManagerSt) AddDeviceCtrl(ctrls DeviceRPCCtrlInterface) {
	thisR.DeviceReqControllersByType[ctrls.GetType()] = ctrls
}

func (this *DeviceRPCManagerSt) RPCReqHandler(c mccommon.ClientToServerHandleResChannel, msg *mccommon.ClientToServerReqSt) error {
	session, col := this.DevicesCollectionManager.GetSesAndCol()
	defer session.Close()

	rpcData := &mccommon.RPCMsg{}
	if err := rpcData.UnmarshalJSON(*msg.Msg); err != nil {
		return this.SendRPCErrorRes(c, msg.Protocol, "", msg.ClientId, 0, err.Error(), 422)
	}

	splitedMethod := strings.Split(rpcData.Method, ".")
	shadowId := splitedMethod[0]
	model := struct{ Type string }{}

	if err := col.Find(&bson.M{
		"shadow.id": shadowId,
	}).One(&model); err != nil {
		return this.SendRPCErrorRes(c, msg.Protocol, rpcData.Method, msg.ClientId, rpcData.Id, err.Error(), 404)
	}

	if this.DeviceReqControllersByType[model.Type] == nil {
		return this.SendRPCErrorRes(c, msg.Protocol, rpcData.Method, msg.ClientId, 0, "No type of device", 404)
	}

	return this.DeviceReqControllersByType[model.Type].HandleReq(rpcData.Method, c, msg, rpcData)
}
