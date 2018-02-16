package mcws

type WsMsgBase struct {}

func (this *WsMsgBase) IsWsMsg() bool {
	return true
}

type WSocketMsgBaseI interface {
	MarshalJSON() ([]byte, error)
	IsWsMsg() bool
}

//easyjson:json
type WsActionMsgBaseSt struct {
	WsMsgBase
	RequestId *string
	Action    string
}

var WsResActionStatuses = map[string]string{
	"success": "success",
	"error": "error",
}

//easyjson:json
type WsResActionMsg struct {
	WsActionMsgBaseSt
	Status string
}

//easyjson:json
type WsResActionSingleErrorMsg struct {
	WsResActionMsg
	Error string
	ErrorCode int
}

func CreateWsResActionSingleErrorMsg(err string, action string, errorCode int, reqId *string) *WsResActionSingleErrorMsg {
	return &WsResActionSingleErrorMsg{
		WsResActionMsg: WsResActionMsg{
			WsActionMsgBaseSt: WsActionMsgBaseSt{
				RequestId: reqId,
				Action: action,
			},
			Status: WsResActionStatuses["error"],
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