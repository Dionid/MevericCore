package mccommon

type DeviceToServerReqHandler func(msg *DeviceToServerReqSt) (res JSONData, sendBack bool, err JSONData)

type DeviceToServerReqSt struct {
	DeviceId  string
	ChannelId string
	Protocol  string
	Msg       *[]byte
}
