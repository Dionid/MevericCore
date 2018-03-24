package mcuserrpcmanager

import (
	"mevericcore/mcuserrpcrouter"
	"mevericcore/mccommunication"
)

type UserRPCManagerBaseSt struct {
	Router *mcuserrpcrouter.UserRPCRouterSt
}

func (thisR *UserRPCManagerBaseSt) RespondRPCErrorRes(c mccommunication.ClientToServerHandleResultChannel, msg *mccommunication.RPCMsg, errMessage string, errCode int) error {
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

func (thisR *UserRPCManagerBaseSt) RespondSuccessResp(c mccommunication.ClientToServerHandleResultChannel, msg *mccommunication.RPCMsg, result *map[string]interface{}) error {
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

func (thisR *UserRPCManagerBaseSt) SendReq(c mccommunication.ClientToServerHandleResultChannel, methodName string, src string, dst string, reqId int, args *map[string]interface{}) error {
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

func (thisR *UserRPCManagerBaseSt) Handle(c mccommunication.ClientToServerHandleResultChannel, msg *mccommunication.ClientToServerRPCReqSt) error {
	return thisR.Router.Handle(c, msg.RPCMsg.Method, msg)
}

func New() *UserRPCManagerBaseSt {
	return &UserRPCManagerBaseSt{
		Router: mcuserrpcrouter.CreateNewDeviceRPCRouter(),
	}
}