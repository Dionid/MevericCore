package mccommon

import "mevericcore/mclibs/mccommunication"

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

func NewShadowUpdateRejectedReqRPC(src string, dst string, methodPrefix string, errMessage string, errCode int) *mccommunication.RPCMsg {
	return &mccommunication.RPCMsg{
		Src: src,
		Dst: dst,
		Method: methodPrefix + ".Device.Shadow.Update.Rejected",
		Args: map[string]interface{}{
			"message": errMessage,
			"code": errCode,
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
