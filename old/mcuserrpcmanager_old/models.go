package mcuserrpcmanager_old

import "mevericcore/mclibs/mcws"

// easyjson:json
type WsAuthRPCReqSt struct {
	mcws.WsRPCMsgBaseSt
	Args struct {
		Login string `json:"login"`
		Password string `json:"password"`
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
