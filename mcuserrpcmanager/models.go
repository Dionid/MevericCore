package mcuserrpcmanager

import "mevericcore/mcws"

// easyjson:json
type WsAuthRPCReqSt struct {
	mcws.WsRPCMsgBaseSt
	Args struct {
		Login string
		Password string
	}
}

// easyjson:json
type WsAuthRPCResSt struct {
	mcws.WsRPCMsgBaseSt
	Result struct{
		Token string
	}
}

// easyjson:json
type WsTokenRPCReqSt struct {
	mcws.WsRPCMsgBaseSt
	Args struct {
		Token string
	}
}
