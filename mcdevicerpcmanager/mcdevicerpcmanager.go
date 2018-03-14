package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
	"strings"
	"gopkg.in/mgo.v2/bson"
)

type ProtocolManagerInterface interface {
	SendJSON(string, mccommon.JSONData) error
}

type DeviceRPCManagerSt struct {
	DeviceResponseServiceSt

	ServerId string

	DevicesCollectionManager mccommon.DevicesCollectionManagerInterface

	DeviceCtrlsByType map[string]DeviceRPCCtrlInterface
}

func CreateDeviceRPCManager(serverId string, devicesColManager mccommon.DevicesCollectionManagerInterface, mqttMan ProtocolManagerInterface) *DeviceRPCManagerSt {
	return &DeviceRPCManagerSt{
		ServerId: serverId,
		DeviceResponseServiceSt: DeviceResponseServiceSt{
			DeviceMQTTManager: mqttMan,
		},
		DevicesCollectionManager: devicesColManager,
		DeviceCtrlsByType: map[string]DeviceRPCCtrlInterface{},
	}
}

func (thisR *DeviceRPCManagerSt) AddDeviceCtrl(deviceType string, ctrls DeviceRPCCtrlInterface) {
	ctrls.SetManagers(thisR.DeviceMQTTManager)
	thisR.DeviceCtrlsByType[deviceType] = ctrls
}

func (this *DeviceRPCManagerSt) RPCReqHandler(msg *mccommon.DeviceToServerReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
	session, col := this.DevicesCollectionManager.GetSesAndCol()
	defer session.Close()

	rpcData := &mccommon.RPCMsg{}
	if err := rpcData.UnmarshalJSON(*msg.Msg); err != nil {
		return this.SendRPCErrorRes(msg.Protocol, msg.DeviceId, 0, err.Error(), 422)
	}

	splitedMethod := strings.Split(rpcData.Method, ".")
	shadowId := splitedMethod[0]
	model := struct{ Type string }{}

	if err := col.Find(&bson.M{
		"shadow.id": shadowId,
	}).One(&model); err != nil {
		return this.SendRPCErrorRes(msg.Protocol, msg.DeviceId, rpcData.Id, err.Error(), 404)
	}

	return this.DeviceCtrlsByType[model.Type].HandleReq(rpcData.Method, msg, rpcData)
}
