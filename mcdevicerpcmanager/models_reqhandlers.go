package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
)

type DeviceResponseServiceSt struct {
	SendRPCMsgToUser mccommon.SendRPCMsgFn
	ServerId         string
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

func (this *DeviceResponseServiceSt) SendRPCShadowDelta(c mccommon.ClientToServerHandleResChannel, protocol string, srcDeviceId string, deviceId string, reqId int, delta *mccommon.ShadowStateDeltaSt) error {
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