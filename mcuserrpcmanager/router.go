package mcuserrpcmanager

import (
	"mevericcore/mccommon"
	"fmt"
	"errors"
)

type (
	ReqSt struct {
		Channel UserRPCManagerHandleResultChannel
		Resource string
		Msg *mccommon.ClientToServerReqSt
		RPCData *mccommon.RPCMsg
		ctx map[string]interface{}
	}

	HandlerFunc func(*ReqSt) error
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

func (r *ReqSt) Get(name string) interface{} {
	return r.ctx[name]
}

func (r *ReqSt) Set(name string, what interface{}) {
	r.ctx[name] = what
}

// ROUTER

type UserRPCRouterSt struct {
	prefix             string
	middlewares         []MiddlewareFunc
	handlersByResource *map[string]HandlerFunc
}

func CreateNewDeviceRPCRouter() *UserRPCRouterSt {
	return &UserRPCRouterSt{
		"",
		[]MiddlewareFunc{},
		&map[string]HandlerFunc{},
	}
}

func (this *UserRPCRouterSt) Group(resource string) *UserRPCRouterSt {
	prefix := ""

	if this.prefix != "" {
		prefix = this.prefix + "."
	}

	return &UserRPCRouterSt{
		prefix:             prefix + resource,
		middlewares: this.middlewares,
		handlersByResource: this.handlersByResource,
	}
}

func (this *UserRPCRouterSt) Use(handler MiddlewareFunc) {
	this.middlewares = append(this.middlewares, handler)
}

func (this *UserRPCRouterSt) AddHandler(resource string, handler HandlerFunc) {
	(*this.handlersByResource)[this.prefix + "." + resource] = handler
}

func (this *UserRPCRouterSt) ChangeAnyHandler(resource string, handler HandlerFunc) {
	(*this.handlersByResource)[resource] = handler
}

func (this *UserRPCRouterSt) Handle(c UserRPCManagerHandleResultChannel, resource string, msg *mccommon.ClientToServerReqSt, rpcData *mccommon.RPCMsg) error {
	defer func() {
		close(c)
		if recover() != nil {
			fmt.Println("Recovered")
			return
		}
	}()

	h := (*this.handlersByResource)[resource]

	if h == nil {
		return errors.New("no handler")
	}

	for i := len(this.middlewares) - 1; i >= 0; i-- {
		h = this.middlewares[i](h)
	}

	return h(&ReqSt{
		c,
		resource,
		msg,
		rpcData,
		map[string]interface{}{},
	})
}
