package mcplantainer

import "mevericcore/mccommunication"

type PlantainerShadowRPCMsgArgsStateSt struct {
	Reported *PlantainerShadowStatePieceSt
	Desired *PlantainerShadowStatePieceSt
}

type PlantainerShadowRPCMsgArgsSt struct {
	State     PlantainerShadowRPCMsgArgsStateSt
	Version   int
}

//easyjson:json
type ShadowUpdateRPCMsgSt struct {
	mccommunication.RPCMsg
	Args PlantainerShadowRPCMsgArgsSt
}
