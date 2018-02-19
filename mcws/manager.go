package mcws

type (
	ReqSt struct {
		Ws *WSocket
		Msg []byte
		ctx map[string]interface{}
	}

	HandlerFunc func(*ReqSt) error
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	ResourcesManagerSt struct {
		prefix             string
		middlewares         []MiddlewareFunc
		handlersByResource *map[string]HandlerFunc
	}
)

func (r *ReqSt) Get(name string) interface{} {
	return r.ctx[name]
}

func (r *ReqSt) Set(name string, what interface{}) {
	r.ctx[name] = what
}

func CreateNewResourcesManager() *ResourcesManagerSt {
	return &ResourcesManagerSt{
		"",
		[]MiddlewareFunc{},
		&map[string]HandlerFunc{},
	}
}

func (this *ResourcesManagerSt) Group(resource string) *ResourcesManagerSt {
	return &ResourcesManagerSt{
		prefix:             this.prefix + "/" + resource,
		middlewares: this.middlewares,
		handlersByResource: this.handlersByResource,
	}
}

func (this *ResourcesManagerSt) Use(handler MiddlewareFunc) {
	this.middlewares = append(this.middlewares, handler)
}

func (this *ResourcesManagerSt) AddHandler(resource string, handler HandlerFunc) {
	(*this.handlersByResource)[this.prefix + "/" + resource] = handler
}

func (this *ResourcesManagerSt) Handle(resource string, ws *WSocket, byteMsg []byte) error {
	h := (*this.handlersByResource)[resource]
	for i := len(this.middlewares) - 1; i >= 0; i-- {
		h = this.middlewares[i](h)
	}
	return h(&ReqSt{
		ws,
		byteMsg,
		map[string]interface{}{},
	})
}
