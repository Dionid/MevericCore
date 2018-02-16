package mcdashboard

import "mevericcore/mcws"

// easyjson:json
type WsTokenActionReqSt struct {
	mcws.WsActionMsgBaseSt
	Login string
	Password string
}

// easyjson:json
type WsTokenActionResSt struct {
	mcws.WsResActionMsg
	token string
}

func CreateAndSendWsTokenActionRes(ws *mcws.WSocket,token string, action string, reqId *string) error {
	res := &WsTokenActionResSt{
		WsResActionMsg: mcws.WsResActionMsg{
			WsActionMsgBaseSt: mcws.WsActionMsgBaseSt{
				RequestId:  reqId,
				Action: action,
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
	mcws.WsActionMsgBaseSt
	Token string
}

// easyjson:json
type WsAuthenticateActionResSt struct {
	mcws.WsResActionMsg
}

func CreateAndSendWsAuthenticateActionRes(ws *mcws.WSocket, action string, reqId *string) error {
	res := &WsAuthenticateActionResSt{
		WsResActionMsg: mcws.WsResActionMsg{
			WsActionMsgBaseSt: mcws.WsActionMsgBaseSt{
				RequestId:  reqId,
				Action: action,
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
