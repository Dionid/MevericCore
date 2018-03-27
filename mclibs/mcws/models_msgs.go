package mcws

import (
	"mevericcore/mclibs/mccommunication"
)

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
	mccommunication.RPCMsg
}

// easyjson:json
type WsAuthRPCReqSt struct {
	WsRPCMsgBaseSt
	Args struct {
		Login string `json:"login"`
		Password string `json:"password"`
	}
}

// easyjson:json
type WsAuthRPCResSt struct {
	WsRPCMsgBaseSt
	Result struct{
		Token string
	}
}

// easyjson:json
type WsTokenRPCReqSt struct {
	WsRPCMsgBaseSt
	Args struct {
		Token string
	}
}