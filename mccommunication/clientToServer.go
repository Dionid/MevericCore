package mccommunication

type JSONData interface {
	MarshalJSON() ([]byte, error)
}

type ClientToServerHandleResult struct {
	Res   JSONData
	Error JSONData
}

type ClientToServerHandleResultChannel chan ClientToServerHandleResult
type ClientToServerReqHandler func(c ClientToServerHandleResultChannel, msg *ClientToServerReqSt) error

type ClientToServerReqSt struct {
	ClientId  string
	Protocol  string
	Resource  *string
	Msg       *[]byte
}

type ClientToServerRPCReqSt struct {
	ClientToServerReqSt
	RPCMsg *RPCMsg
}
