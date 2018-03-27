package mcrpcmanager


import (
	"mevericcore/mclibs/mccommunication"
	"mevericcore/mclibs/mcrpcrouter"
)

type RPCManagerSt struct {
	Router *mcrpcrouter.RPCRouterSt
}

func (thisR *RPCManagerSt) RespondRPCErrorRes(c mccommunication.ClientToServerHandleResultChannel, msg *mccommunication.RPCMsg, errMessage string, errCode int) error {
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

func (thisR *RPCManagerSt) RespondSuccessResp(c mccommunication.ClientToServerHandleResultChannel, msg *mccommunication.RPCMsg, result *map[string]interface{}) error {
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

func (thisR *RPCManagerSt) SendReq(c mccommunication.ClientToServerHandleResultChannel, methodName string, src string, dst string, reqId int, args *map[string]interface{}) error {
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

func (thisR *RPCManagerSt) Handle(c mccommunication.ClientToServerHandleResultChannel, rpcData *mccommunication.RPCMsg, msg *[]byte) error {
	return thisR.Router.Handle(c, rpcData.Method, mcrpcrouter.NewRPCReqSt(c, rpcData, msg))
}

func New() *RPCManagerSt {
	return &RPCManagerSt{
		mcrpcrouter.New(),
	}
}
