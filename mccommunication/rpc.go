package mccommunication

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

type SendRPCMsgFn func(msg *RPCMsg) error
