package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
)

type DeviceResponseServiceSt struct {
	DeviceHTTPManager ProtocolManagerInterface
	DeviceMQTTManager ProtocolManagerInterface
	DeviceWSManager ProtocolManagerInterface

	SendToUser func(msg *mccommon.RPCMsg) error

	ServerId string
}

func (this *DeviceResponseServiceSt) SendRPCErrorRes(c mccommon.ClientToServerHandleResChannel, protocol string, methodName string, srcDeviceId string, reqId int, errMessage string, errCode int) error {
	data := mccommon.RPCMsg{
		Src: this.ServerId,
		Dst: srcDeviceId,
		Method: methodName,
		Id: reqId,
		Error: &map[string]interface{}{
			"message": errMessage,
			"code": errCode,
		},
	}

	c <- mccommon.ClientToServerHandleRes{
		nil,
		data,
	}

	return nil
}

func (this *DeviceResponseServiceSt) SendRPCSuccessRes(c mccommon.ClientToServerHandleResChannel, protocol string, methodName string, srcDeviceId string, reqId int, result *map[string]interface{}) error {
	data := mccommon.RPCMsg{
		Src: this.ServerId,
		Dst: srcDeviceId,
		Method: methodName,
		Id: reqId,
		Result: result,
	}

	c <- mccommon.ClientToServerHandleRes{
		data,
		nil,
	}

	return nil
}

func (this *DeviceResponseServiceSt) PublishDelta(c mccommon.ClientToServerHandleResChannel, protocol string, srcDeviceId string, deviceId string, reqId int, delta *mccommon.ShadowStateDeltaSt) error {
	//if protocol == "MQTT" {
	//	this.DeviceMQTTManager.SendJSON(srcDeviceId+"/rpc", mccommon.RPCMsg{
	//		Src: this.ServerId,
	//		Dst: srcDeviceId,
	//		Id: reqId,
	//		Method: deviceId+".Shadow.Delta",
	//		Args: delta,
	//	})
	//}
	//if protocol == "WS" {
	//	this.DeviceWSManager.SendJSON(srcDeviceId, mccommon.RPCMsg{
	//		Src: this.ServerId,
	//		Dst: srcDeviceId,
	//		Id: reqId,
	//		Method: deviceId+".Shadow.Delta",
	//		Args: delta,
	//	})
	//}

	data := mccommon.RPCMsg{
		Src: this.ServerId,
		Dst: srcDeviceId,
		Id: reqId,
		Method: deviceId+".Shadow.Delta",
		Args: delta,
	}

	c <- mccommon.ClientToServerHandleRes{
		nil,
		data,
	}

	return nil
}