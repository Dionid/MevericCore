package mcdashboard

import "mevericcore/mcws"

// easyjson:json
type WsTokenActionMsgSt struct {
	mcws.WsActionMsgBaseSt
	Login string
	Password string
}


// easyjson:json
type WsAuthenticateActionMsgSt struct {
	mcws.WsActionMsgBaseSt
	Token string
}
