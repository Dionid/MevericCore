package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
	"strings"
)

type (
	ReqSt struct {
		Resource string
		Msg *mccommon.DeviceToServerReqSt
		RPCData *mccommon.RPCMsg
		DeviceId string
		ctx map[string]interface{}
	}

	HandlerFunc func(*ReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData)
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

func (r *ReqSt) Get(name string) interface{} {
	return r.ctx[name]
}

func (r *ReqSt) Set(name string, what interface{}) {
	r.ctx[name] = what
}

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
	return &DeviceRPCRouterSt{
		prefix:             this.prefix + "." + resource,
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

func (this *DeviceRPCRouterSt) Handle(resource string, msg *mccommon.DeviceToServerReqSt, rpcData *mccommon.RPCMsg) (mccommon.JSONData, bool, mccommon.JSONData) {
	splitedRes := strings.Split(resource, ".")
	res := strings.Join(splitedRes[1:], ".")

	h := (*this.handlersByResource)[res]
	for i := len(this.middlewares) - 1; i >= 0; i-- {
		h = this.middlewares[i](h)
	}

	return h(&ReqSt{
		resource,
		msg,
		rpcData,
		splitedRes[0],
		map[string]interface{}{},
	})
}

//type DeviceSpecificRouterHandlersSt struct{
//	BeforeHandlerExec HandlerFunc
//	AfterHandlerExec HandlerFunc
//}

//type DeviceSpecificRouter struct {
//	Handlers map[string]*DeviceSpecificRouterHandlersSt
//}
//
//func (this *DeviceSpecificRouter) AddHandler(resource string, beforeH HandlerFunc, afterH HandlerFunc) {
//	this.Handlers[resource] = &DeviceSpecificRouterHandlersSt{
//		beforeH,
//		afterH,
//	}
//}
//
//func (this *DeviceSpecificRouter) BeforeHandlerExec(req *ReqSt) {
//	this.Handlers[req.Resource].BeforeHandlerExec(req)
//}
//
//func (this *DeviceSpecificRouter) AfterHandlerExec(req *ReqSt) {
//	this.Handlers[req.Resource].AfterHandlerExec(req)
//}
//
//func (this *DeviceSpecificRouter) ExecDeviceSpecificRouterMiddleware(h HandlerFunc) HandlerFunc {
//	return func(req *ReqSt) (mccommon.JSONData, bool, mccommon.JSONData) {
//		this.BeforeHandlerExec(req)
//
//		res, s, e := h(req)
//
//		this.AfterHandlerExec(req)
//
//		return res, s, e
//	}
//}