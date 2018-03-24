package mcdevicerpcmanager

import (
	"mevericcore/mccommunication"
	"mevericcore/mcdevicerpcrouter"
)

type DeviceRPCManagerSt struct {
	Router *mcdevicerpcrouter.DeviceRPCRouterSt
}

func (thisR *DeviceRPCManagerSt) RespondRPCErrorRes(c mccommunication.ClientToServerHandleResultChannel, msg *mccommunication.RPCMsg, errMessage string, errCode int) error {
	data := &mccommunication.RPCMsg{
		Method: msg.Method,
		Id: msg.Id,
		Src: msg.Dst,
		Dst: msg.Src,
		Error: &map[string]interface{}{
			"message": errMessage,
			"code": errCode,
		},
	}
	c <- mccommunication.ClientToServerHandleResult{
		nil,
		data,
	}
	return nil
}

func (thisR *DeviceRPCManagerSt) RespondSuccessResp(c mccommunication.ClientToServerHandleResultChannel, msg *mccommunication.RPCMsg, result *map[string]interface{}) error {
	data := &mccommunication.RPCMsg{
		Method: msg.Method,
		Id: msg.Id,
		Src: msg.Dst,
		Dst: msg.Src,
		Result: result,
	}
	c <- mccommunication.ClientToServerHandleResult{
		data,
		nil,
	}
	return nil
}

func (thisR *DeviceRPCManagerSt) SendReq(c mccommunication.ClientToServerHandleResultChannel, methodName string, src string, dst string, reqId int, args *map[string]interface{}) error {
	data := &mccommunication.RPCMsg{
		Method: methodName,
		Id: reqId,
		Src: src,
		Dst: dst,
		Args: args,
	}
	c <- mccommunication.ClientToServerHandleResult{
		data,
		nil,
	}
	return nil
}

func (thisR *DeviceRPCManagerSt) Handle(c mccommunication.ClientToServerHandleResultChannel, msg *mccommunication.ClientToServerRPCReqSt) error {
	return thisR.Router.Handle(c, msg.RPCMsg.Method, msg)
}

func New() *DeviceRPCManagerSt {
	return &DeviceRPCManagerSt{
		mcdevicerpcrouter.New(),
	}
}