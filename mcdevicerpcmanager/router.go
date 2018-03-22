package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
	"strings"
	"fmt"
	"errors"
)

// REQ

type (
	ReqSt struct {
		Channel mccommon.ClientToServerHandleResultChannel
		Resource string
		Msg *mccommon.ClientToServerReqSt
		RPCData *mccommon.RPCMsg
		DeviceId string
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

type DeviceRPCRouterSt struct {
	prefix             string
	middlewares         []MiddlewareFunc
	handlersByResource *map[string]HandlerFunc
}

func CreateNewDeviceRPCRouter() *DeviceRPCRouterSt {
	return &DeviceRPCRouterSt{
		"",
		[]MiddlewareFunc{},
		&map[string]HandlerFunc{},
	}
}

func (this *DeviceRPCRouterSt) Group(resource string) *DeviceRPCRouterSt {
	prefix := ""

	if this.prefix != "" {
		prefix = this.prefix + "."
	}

	return &DeviceRPCRouterSt{
		prefix:             prefix + resource,
		middlewares: this.middlewares,
		handlersByResource: this.handlersByResource,
	}
}

func (this *DeviceRPCRouterSt) Use(handler MiddlewareFunc) {
	this.middlewares = append(this.middlewares, handler)
}

func (this *DeviceRPCRouterSt) AddHandler(resource string, handler HandlerFunc) {
	(*this.handlersByResource)[this.prefix + "." + resource] = handler
}

func (this *DeviceRPCRouterSt) ChangeAnyHandler(resource string, handler HandlerFunc) {
	(*this.handlersByResource)[resource] = handler
}

func (this *DeviceRPCRouterSt) Handle(resource string, c mccommon.ClientToServerHandleResultChannel, msg *mccommon.ClientToServerReqSt, rpcData *mccommon.RPCMsg) error {
	defer func() {
		if recover() != nil {
			fmt.Println("Recovered")
			return
		}
	}()

	splitedRes := strings.Split(resource, ".")
	res := strings.Join(splitedRes[1:], ".")

	h := (*this.handlersByResource)[res]

	if h == nil {
		return errors.New("no handler")
	}

	for i := len(this.middlewares) - 1; i >= 0; i-- {
		h = this.middlewares[i](h)
	}

	return h(&ReqSt{
		c,
		res,
		msg,
		rpcData,
		splitedRes[0],
		map[string]interface{}{},
	})
}