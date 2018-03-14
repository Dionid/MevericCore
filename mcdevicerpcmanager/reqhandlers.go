package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
)

type DeviceResponseServiceSt struct {
	DeviceHTTPManager ProtocolManagerInterface
	DeviceMQTTManager ProtocolManagerInterface
	DeviceWSManager ProtocolManagerInterface

	ServerId string
}

func (this *DeviceResponseServiceSt) SendRPCErrorRes(protocol string, srcDeviceId string, reqId int, errMessage string, errCode int) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
	return nil, true, mccommon.RPCMsg{
		Src: this.ServerId,
		Dst: srcDeviceId,
		Id: reqId,
		Error: &map[string]interface{}{
			"message": errMessage,
			"code": errCode,
		},
	}
}

func (this *DeviceResponseServiceSt) PublishDelta(protocol string, srcDeviceId string, deviceId string, reqId int, delta *mccommon.ShadowStateDeltaSt) {
	if protocol == "MQTT" {
		this.DeviceMQTTManager.SendJSON(srcDeviceId+"/rpc", mccommon.RPCMsg{
			Src: this.ServerId,
			Dst: srcDeviceId,
			Id: reqId,
			Method: deviceId+".Shadow.Delta",
			Args: delta,
		})
	}
	if protocol == "WS" {
		this.DeviceWSManager.SendJSON(srcDeviceId, mccommon.RPCMsg{
			Src: this.ServerId,
			Dst: srcDeviceId,
			Id: reqId,
			Method: deviceId+".Shadow.Delta",
			Args: delta,
		})
	}
}

//func (this *DeviceRPCManagerSt) SendRPCErrorRes(protocol string, srcDeviceId string, reqId int, errMessage string, errCode int) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
//	return nil, true, mccommon.RPCMsg{
//		Src: this.ServerId,
//		Dst: srcDeviceId,
//		Id: reqId,
//		Error: &map[string]interface{}{
//			"message": errMessage,
//			"code": errCode,
//		},
//	}
//}
//
//func (this *DeviceRPCManagerSt) PublishDelta(protocol string, srcDeviceId string, deviceId string, reqId int, delta *mccommon.ShadowStateDeltaSt) {
//	if protocol == "MQTT" {
//		this.DeviceMQTTManager.SendJSON(srcDeviceId+"/rpc", mccommon.RPCMsg{
//			Src: this.ServerId,
//			Dst: srcDeviceId,
//			Id: reqId,
//			Method: deviceId+".Shadow.Delta",
//			Args: delta,
//		})
//	}
//	if protocol == "WS" {
//		this.DeviceWSManager.SendJSON(srcDeviceId, mccommon.RPCMsg{
//			Src: this.ServerId,
//			Dst: srcDeviceId,
//			Id: reqId,
//			Method: deviceId+".Shadow.Delta",
//			Args: delta,
//		})
//	}
//}

//func (this *DeviceRPCManagerSt) RPCReqHandler(msg *mccommon.DeviceToServerReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
//	rpcData := &mccommon.RPCMsg{}
//	if err := rpcData.UnmarshalJSON(*msg.Msg); err != nil {
//		return this.SendRPCErrorRes(msg.Protocol, msg.DeviceId, 0, err.Error(), 422)
//	}
//
//	splitedMethod := strings.Split(rpcData.Method, ".")
//	shadowId := splitedMethod[0]
//	model := struct{ Type string }{}
//
//	if err := this.DevicesCollectionManager.Find(&bson.M{
//		"shadow.id": shadowId,
//	}).One(&model); err != nil {
//		return this.SendRPCErrorRes(msg.Protocol, msg.DeviceId, rpcData.Id, err.Error(), 404)
//	}
//
//	//device.SetIsActivated(colManager)
//
//	return this.Router.Handle(rpcData.Method, msg, rpcData)
//}
