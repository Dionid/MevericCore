package mcuserrpcrouter

import (
	"fmt"
	"errors"
	"mevericcore/mccommunication"
)

type UserRPCRouterSt struct {
	prefix             string
	middlewares         []mccommunication.RPCMiddlewareFunc
	handlersByResource *map[string]mccommunication.RPCHandlerFunc
}

func CreateNewDeviceRPCRouter() *UserRPCRouterSt {
	return &UserRPCRouterSt{
		"",
		[]mccommunication.RPCMiddlewareFunc{},
		&map[string]mccommunication.RPCHandlerFunc{},
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

func (this *UserRPCRouterSt) AddHandler(resource string, handler mccommunication.RPCHandlerFunc) {
	prefix := resource

	if this.prefix != ""{
		prefix = this.prefix + "." + resource
	}

	(*this.handlersByResource)[prefix] = handler
}

func (this *UserRPCRouterSt) ChangeAnyHandler(resource string, handler mccommunication.RPCHandlerFunc) {
	(*this.handlersByResource)[resource] = handler
}

func (this *UserRPCRouterSt) Handle(c mccommunication.ClientToServerHandleResultChannel, resource string, msg *mccommunication.ClientToServerRPCReqSt) error {
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

	return h(mccommunication.NewRPCReqSt(c, msg))
}
