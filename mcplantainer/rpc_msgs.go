package mcplantainer

import (
	"mevericcore/mclibs/mccommunication"
	"mevericcore/mclibs/mccommon"
)

func NewShadowUpdateDeltaReqRPC(dst string, shadow *PlantainerShadowSt) *mccommunication.RPCMsg {
	return mccommon.NewShadowUpdateDeltaReqRPC(PlantainerServerId, dst, "Plantainer", shadow.State.Delta, shadow.Metadata.Version)
}

func NewShadowUpdateAcceptedReqRPC(dst string, shadow *PlantainerShadowSt) *mccommunication.RPCMsg {
	return mccommon.NewShadowUpdateAcceptedReqRPC(PlantainerServerId, dst, "Plantainer", shadow.State, shadow.Metadata.Version)
}