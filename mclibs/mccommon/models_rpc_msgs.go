package mccommon

type ClientToServerHandleResult struct {
	Res   JSONData
	Error JSONData
}

type ClientToServerHandleResultChannel chan ClientToServerHandleResult

type ClientToServerReqHandler func(c ClientToServerHandleResultChannel, msg *ClientToServerReqSt) error

type ClientToServerReqSt struct {
	ClientId  string
	ChannelId string
	Protocol  string
	Resource *string
	Msg       *[]byte
}

type SendRPCMsgFn func(msg *RPCMsg) error

//easyjson:json
type RPCMsg struct {
	Src string
	Dst string
	Method string
	Id  int
	Args interface{}
	Error *map[string]interface{}
	Result *map[string]interface{}
}

//easyjson:json
type RPCWithShadowUpdateMsg struct {
	RPCMsg
	Args DeviceShadowUpdateMsg
}

