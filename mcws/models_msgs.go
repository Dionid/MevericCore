package mcws

import "mevericcore/mccommon"

type WsMsgBase struct {}

func (this *WsMsgBase) IsWsMsg() bool {
	return true
}

type WSocketMsgBaseI interface {
	MarshalJSON() ([]byte, error)
	IsWsMsg() bool
}

//easyjson:json
type WsRPCMsgBaseSt struct {
	WsMsgBase
	mccommon.RPCMsg
}

type WsResActionStatusesSt struct {
	Success string
	Error string
}

var WsResActionStatuses = WsResActionStatusesSt{
	"success",
	"error",
}

//easyjson:json
type WsResActionMsg struct {
	WsRPCMsgBaseSt
	Status string
}

//easyjson:json
type WsResActionSingleErrorMsg struct {
	WsResActionMsg
	Error string
	ErrorCode int
}

func CreateWsResActionSingleErrorMsg(err string, action string, errorCode int, reqId int) *WsResActionSingleErrorMsg {
	return &WsResActionSingleErrorMsg{
		WsResActionMsg: WsResActionMsg{
			WsRPCMsgBaseSt: WsRPCMsgBaseSt{
				RPCMsg: mccommon.RPCMsg{
					Method: action,
					Id:     reqId,
				},
			},
			Status: WsResActionStatuses.Error,
		},
		Error: err,
		ErrorCode: errorCode,
	}
}

//easyjson:json
type WsResActionArrErrorMsg struct {
	WsResActionMsg
	Errors []string
}