package mcdashboard_old

import (
	"mevericcore/mclibs/mcws"
	"mevericcore/mclibs/mccommon"
)

// easyjson:json
type WsTokenActionReqSt struct {
	mcws.WsRPCMsgBaseSt
	Login string
	Password string
}

// easyjson:json
type WsTokenActionResSt struct {
	mcws.WsResActionMsg
	token string
}

func CreateAndSendWsTokenActionRes(ws *mcws.WSocket,token string, action string, reqId int) error {
	res := &WsTokenActionResSt{
		WsResActionMsg: mcws.WsResActionMsg{
			WsRPCMsgBaseSt: mcws.WsRPCMsgBaseSt{
				WsMsgBase: nil,
				RPCMsg: mccommon.RPCMsg{
					Id:     reqId,
					Method: action,
				},
			},
			Status: mcws.WsResActionStatuses.Success,
		},
		token: token,
	}
	if bData, err := res.MarshalJSON(); err != nil {
		return err
	} else {
		ws.SendMsg(bData)
	}

	return nil
}

// easyjson:json
type WsAuthenticateActionReqSt struct {
	mcws.WsRPCMsgBaseSt
	Token string
}

// easyjson:json
type WsAuthenticateActionResSt struct {
	mcws.WsResActionMsg
}

func CreateAndSendWsAuthenticateActionRes(ws *mcws.WSocket, action string, reqId *string) error {
	res := &WsAuthenticateActionResSt{
		WsResActionMsg: mcws.WsResActionMsg{
			WsRPCMsgBaseSt: mcws.WsRPCMsgBaseSt{
				Id:     reqId,
				Method: action,
			},
			Status: mcws.WsResActionStatuses.Success,
		},
	}
	if bData, err := res.MarshalJSON(); err != nil {
		return err
	} else {
		ws.SendMsg(bData)
	}

	return nil
}
