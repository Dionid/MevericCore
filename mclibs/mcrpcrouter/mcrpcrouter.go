package mcrpcrouter

import (
	"fmt"
	"errors"
	"mevericcore/mclibs/mccommunication"
	"strings"
)

type (
	RPCReqSt struct {
		mccommunication.ReqSt
		Msg *[]byte
		RPCData *mccommunication.RPCMsg
	}

	RPCHandlerFunc func(*RPCReqSt) error
	RPCMiddlewareFunc func(RPCHandlerFunc) RPCHandlerFunc
)

func NewRPCReqSt(c mccommunication.ClientToServerHandleResultChannel, rpcData *mccommunication.RPCMsg, msg *[]byte) *RPCReqSt {
	return &RPCReqSt{
		ReqSt: *mccommunication.NewReqSt(c),
		Msg: msg,
		RPCData: rpcData,
	}
}

type handlersAndChildSt struct {
	handler *RPCHandlerFunc
	children *map[string]*handlersAndChildSt
}

type RPCRouterSt struct {
	prefix             string
	middlewares         []RPCMiddlewareFunc
	handlersByResource *map[string]*handlersAndChildSt
}

func New() *RPCRouterSt {
	return &RPCRouterSt{
		"",
		[]RPCMiddlewareFunc{},
		&map[string]*handlersAndChildSt{},
	}
}

func (this *RPCRouterSt) Group(resource string) *RPCRouterSt {
	prefix := ""

	if this.prefix != "" {
		prefix = this.prefix + "."
	}

	return &RPCRouterSt{
		prefix:             prefix + resource,
		middlewares: this.middlewares,
		handlersByResource: this.handlersByResource,
	}
}

func (this *RPCRouterSt) Use(handler RPCMiddlewareFunc) {
	this.middlewares = append(this.middlewares, handler)
}

func (this *RPCRouterSt) AddHandler(res string, handler RPCHandlerFunc) {
	fullResource := res

	if this.prefix != ""{
		fullResource = this.prefix + "." + res
	}

	fullResourceSep := strings.Split(fullResource, ".")

	handlersAndChildrenMap := this.handlersByResource

	for index, resourceName := range fullResourceSep {
		h := (*handlersAndChildrenMap)[resourceName]
		if h == nil {
			if len(fullResourceSep) - 1 == index {
				(*handlersAndChildrenMap)[resourceName] = &handlersAndChildSt{
					children: &map[string]*handlersAndChildSt{},
					handler: &handler,
				}
			} else {
				(*handlersAndChildrenMap)[resourceName] = &handlersAndChildSt{
					children: &map[string]*handlersAndChildSt{},
					handler: nil,
				}
				handlersAndChildrenMap = (*handlersAndChildrenMap)[resourceName].children
			}
		} else {
			if len(fullResourceSep) - 1 == index {
				h.handler = &handler
			} else {
				handlersAndChildrenMap = h.children
			}
		}
	}
}

//func (this *RPCRouterSt) ChangeAnyHandler(resource string, handler mccommunication.RPCHandlerFunc) {
//	(*this.handlersByResource)[resource] = handler
//}

func (this *RPCRouterSt) Handle(c mccommunication.ClientToServerHandleResultChannel, resource string, msg *RPCReqSt) error {
	defer func() {
		close(c)
		if recover() != nil {
			fmt.Println("RPCRouterSt.Handle recovered")
			return
		}
	}()

	resourceSep := strings.Split(resource, ".")

	handlersAndChildrenMap := this.handlersByResource

	for index, resourceName := range resourceSep {
		h := (*handlersAndChildrenMap)[resourceName]

		if h == nil {
			if (*handlersAndChildrenMap)["*"] != nil {
				return (*(*handlersAndChildrenMap)["*"].handler)(msg)
			}
			if (*handlersAndChildrenMap)["#"] != nil {
				if len(resourceSep) - 1 == index {
					if (*handlersAndChildrenMap)["#"].handler == nil {
						return errors.New("no handler on resource: " + resource)
					} else {
						return (*(*handlersAndChildrenMap)["#"].handler)(msg)
					}
				} else {
					handlersAndChildrenMap = (*handlersAndChildrenMap)["#"].children
					continue
				}
			}
			return errors.New("no resource: " + resource)
		}

		if len(resourceSep) - 1 == index {
			if h.handler == nil {
				return errors.New("no handler on resource: " + resource)
			} else {
				return (*h.handler)(msg)
			}
		} else {
			if h.children == nil {
				return errors.New("no subresources on resource: " + resource)
			} else {
				handlersAndChildrenMap = h.children
			}
		}
	}

	return nil
}