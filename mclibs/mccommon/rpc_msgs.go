package mccommon

import "mevericcore/mclibs/mccommunication"

////easyjson:json
//type ShadowUpdateAcceptedReqArgsStateSt struct {
//	Reported *map[string]interface{}
//	Desired  *map[string]interface{}
//}

////easyjson:json
//type ShadowUpdateAcceptedReqArgsSt struct {
//	State interface{}
//	Version int
//}

func NewShadowUpdateAcceptedReqRPC(src string, dst string, methodPrefix string, state interface{}, version int) *mccommunication.RPCMsg {
	return &mccommunication.RPCMsg{
		Src: src,
		Dst: dst,
		Method: methodPrefix + ".Device.Shadow.Update.Accepted",
		Args: map[string]interface{}{
			"state": state,
			"version": version,
		},
	}
}

func NewShadowUpdateDeltaReqRPC(src string, dst string, methodPrefix string, delta interface{}, version int) *mccommunication.RPCMsg {
	return &mccommunication.RPCMsg{
		Src: src,
		Dst: dst,
		Method: methodPrefix + ".Device.Shadow.Delta",
		Args: map[string]interface{}{
			"state": delta,
			"version": version,
		},
	}
}
