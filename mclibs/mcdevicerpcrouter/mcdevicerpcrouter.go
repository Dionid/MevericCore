package mcdevicerpcrouter

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

type DeviceRPCRouterSt struct {
	prefix             string
	middlewares         []mccommunication.RPCMiddlewareFunc
	handlersByResource *map[string]*handlersAndChildSt
}

func New() *DeviceRPCRouterSt {
	return &DeviceRPCRouterSt{
		"",
		[]mccommunication.RPCMiddlewareFunc{},
		&map[string]*handlersAndChildSt{},
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

func (this *DeviceRPCRouterSt) Use(handler mccommunication.RPCMiddlewareFunc) {
	this.middlewares = append(this.middlewares, handler)
}

func (this *DeviceRPCRouterSt) AddHandler(res string, handler mccommunication.RPCHandlerFunc) {
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

//func (this *DeviceRPCRouterSt) ChangeAnyHandler(resource string, handler mccommunication.RPCHandlerFunc) {
//	(*this.handlersByResource)[resource] = handler
//}

func (this *DeviceRPCRouterSt) Handle(c mccommunication.ClientToServerHandleResultChannel, resource string, msg *mccommunication.ClientToServerRPCReqSt) error {
	defer func() {
		close(c)
		if recover() != nil {
			fmt.Println("DeviceRPCRouterSt.Handle Recovered")
			return
		}
	}()

	resourceSep := strings.Split(resource, ".")

	handlersAndChildrenMap := this.handlersByResource

	for index, resourceName := range resourceSep {
		h := (*handlersAndChildrenMap)[resourceName]

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
}

