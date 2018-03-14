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
	ServerId string

	DeviceResponseServiceSt
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
	}
}

func (thisR *DeviceRPCManagerSt) AddDeviceCtrl(deviceType string, ctrls DeviceRPCCtrlInterface) {
	ctrls.SetManagers(thisR.DeviceMQTTManager)
	thisR.DeviceCtrlsByType[deviceType] = ctrls
}

func (this *DeviceRPCManagerSt) RPCReqHandler(msg *mccommon.DeviceToServerReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
	rpcData := &mccommon.RPCMsg{}
	if err := rpcData.UnmarshalJSON(*msg.Msg); err != nil {
		return this.SendRPCErrorRes(msg.Protocol, msg.DeviceId, 0, err.Error(), 422)
	}

	splitedMethod := strings.Split(rpcData.Method, ".")
	shadowId := splitedMethod[0]
	model := struct{ Type string }{}

	if err := this.DevicesCollectionManager.Find(&bson.M{
		"shadow.id": shadowId,
	}).One(&model); err != nil {
		return this.SendRPCErrorRes(msg.Protocol, msg.DeviceId, rpcData.Id, err.Error(), 404)
	}

	return this.DeviceCtrlsByType[model.Type].HandleReq(rpcData.Method, msg, rpcData)
}

//func (thisR *DeviceRPCManagerSt) initShadowRoutes() {
//	shadowG := thisR.Router.Group("Shadow")
//	shadowG.AddHandler("Get", func(req *ReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
//		state := req.Device.GetShadow().GetState()
//		state.FillDelta()
//
//		if len(state.Delta.State) != 0 {
//			thisR.PublishDelta(req.Msg.Protocol, req.Msg.DeviceId, req.DeviceId, req.RPCData.Id, state.Delta)
//		}
//
//		return mccommon.RPCMsg{
//			Src: thisR.ServerId,
//			Dst: req.Msg.DeviceId,
//			Id: req.RPCData.Id,
//			Result: &map[string]interface{}{
//				"state": state,
//			},
//		}, true, nil
//	})
//}

func (this *DeviceRPCManagerSt) Init() {
	//this.initShadowRoutes()
}

