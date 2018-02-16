package mcws

type WsMsgBase struct {
	RequestId *string
	Action    string
}

func (this *WsMsgBase) IsWsMsg() bool {
	return true
}

type WSocketMsgBaseI interface {
	MarshalJSON() ([]byte, error)
	IsWsMsg() bool
}
