package mcuserrpcrouter

import (
	"fmt"
	"errors"
	"mevericcore/mclibs/mccommunication"
	"strings"
)

type handlersAndChildSt struct {
	handler *mccommunication.RPCHandlerFunc
	children *map[string]*handlersAndChildSt
}

type UserRPCRouterSt struct {
	prefix             string
	middlewares         []mccommunication.RPCMiddlewareFunc
	handlersByResource *map[string]*handlersAndChildSt
}

func CreateNewDeviceRPCRouter() *UserRPCRouterSt {
	return &UserRPCRouterSt{
		"",
		[]mccommunication.RPCMiddlewareFunc{},
		&map[string]*handlersAndChildSt{},
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

func (this *UserRPCRouterSt) Use(handler mccommunication.RPCMiddlewareFunc) {
	this.middlewares = append(this.middlewares, handler)
}

func (this *UserRPCRouterSt) AddHandler(res string, handler mccommunication.RPCHandlerFunc) {
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

//func (this *UserRPCRouterSt) ChangeAnyHandler(resource string, handler mccommunication.RPCHandlerFunc) {
//	(*this.handlersByResource)[resource] = handler
//}

func (this *UserRPCRouterSt) Handle(c mccommunication.ClientToServerHandleResultChannel, resource string, msg *mccommunication.ClientToServerRPCReqSt) error {
	defer func() {
		close(c)
		if recover() != nil {
			fmt.Println("UserRPCRouterSt.Handle recovered")
			return
		}
	}()

	resourceSep := strings.Split(resource, ".")

	handlersAndChildrenMap := this.handlersByResource

	for index, resourceName := range resourceSep {
		h := (*handlersAndChildrenMap)[resourceName]

		//if h == nil && index == 0 {
		//	return errors.New("no resource: " + resource)
		//}

		if h == nil {
			if (*handlersAndChildrenMap)["*"] != nil {
				return (*(*handlersAndChildrenMap)["*"].handler)(mccommunication.NewRPCReqSt(c, msg))
			}
			if (*handlersAndChildrenMap)["#"] != nil {
				if len(resourceSep) - 1 == index {
					if (*handlersAndChildrenMap)["#"].handler == nil {
						return errors.New("no handler on resource: " + resource)
					} else {
						return (*(*handlersAndChildrenMap)["#"].handler)(mccommunication.NewRPCReqSt(c, msg))
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
				return (*h.handler)(mccommunication.NewRPCReqSt(c, msg))
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

	//if h == nil {
	//	return errors.New("no handler")
	//}
	//
	////for i := len(this.middlewares) - 1; i >= 0; i-- {
	////	h = this.middlewares[i](h)
	////}
	//
	//return h(mccommunication.NewRPCReqSt(c, msg))
}
