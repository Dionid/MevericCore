package mcinnerrpc

import (
	"github.com/nats-io/go-nats"
)

type Msg struct {
	Data []byte
}

type MsgHandler func(*Msg)

type InnerRPCServiceSt struct {
	OriginalService  *nats.Conn
}

func NewInnerRPCService() *InnerRPCServiceSt {
	nc, _ := nats.Connect(nats.DefaultURL)
	return &InnerRPCServiceSt{
		OriginalService: nc,
	}
}

func (this *InnerRPCServiceSt) Subscribe(subj string, cb MsgHandler) (interface{}, error) {
	return this.OriginalService.Subscribe(subj, func(msg *nats.Msg) {
		newMsg := &Msg{
			msg.Data,
		}
		cb(newMsg)
	})
}

func (this *InnerRPCServiceSt) Publish(subj string, data []byte) error {
	return this.OriginalService.Publish(subj, data)
}

type InnerRPCServiceInterface interface {
	Publish(string, []byte) error
	Subscribe(subj string, cb MsgHandler) (interface{}, error)
}

type InnerRPCManSt struct {
	Service InnerRPCServiceInterface
}

func New() *InnerRPCManSt {
	return &InnerRPCManSt{}
}

func (this *InnerRPCManSt) Init() {
	this.Service = NewInnerRPCService()
}

//func (this *InnerRPCManSt) SendRPCMsgToUser(msg *mccommon.RPCMsg) error {
//	if bData, err := msg.MarshalJSON(); err != nil {
//		return err
//	} else {
//		this.Service.Publish("User.RPC.Send", bData)
//	}
//
//	return nil
//}
