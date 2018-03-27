package mccommunication

type (
	ReqSt struct {
		Channel ClientToServerHandleResultChannel
		Msg *ClientToServerReqSt
		ctx map[string]interface{}
	}

	HandlerFunc func(*ReqSt) error
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	RPCReqSt struct {
		ReqSt
		Msg *ClientToServerRPCReqSt
	}

	RPCHandlerFunc func(*RPCReqSt) error
	RPCMiddlewareFunc func(RPCHandlerFunc) RPCHandlerFunc
)

func NewReqSt(c ClientToServerHandleResultChannel) *ReqSt {
	return &ReqSt{
		Channel: c,
		ctx: map[string]interface{}{},
	}
}

func NewRPCReqSt(c ClientToServerHandleResultChannel, msg *ClientToServerRPCReqSt) *RPCReqSt {
	return &RPCReqSt{
		ReqSt: ReqSt{
			Channel: c,
			ctx: map[string]interface{}{},
		},
		Msg: msg,
	}
}

func (r *ReqSt) Get(name string) interface{} {
	return r.ctx[name]
}

func (r *ReqSt) Set(name string, what interface{}) {
	r.ctx[name] = what
}
