package mccommon

type DeviceToServerReqHandler func(msg *DeviceToServerReqSt) (res JSONData, sendBack bool, err JSONData)

type DeviceToServerReqSt struct {
	DeviceId  string
	ChannelId string
	Protocol  string
	Msg       *[]byte
}

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