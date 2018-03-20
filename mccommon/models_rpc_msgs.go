package mccommon

type ClientToServerHandleRes struct {
	Resp JSONData
	Error JSONData
}

type ClientToServerHandleResChannel chan ClientToServerHandleRes

type ClientToServerReqHandler func(c ClientToServerHandleResChannel, msg *ClientToServerReqSt) error
//type ClientToServerRPCReqHandler func(c *ClientToServerHandleResChannel, msg *ClientToServerRPCReqSt) error


type ClientToServerReqSt struct {
	ClientId  string
	ChannelId string
	Protocol  string
	Resource *string
	Msg       *[]byte
}

//type ClientToServerRPCReqSt struct {
//	ClientToServerReqSt
//	RPCMsg *RPCMsg
//}

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