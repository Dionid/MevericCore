package mcplantainer

import (
	"mevericcore/mclibs/mccommunication"
	"fmt"
)

// ToDo: separate this to another file

type LightModuleStateDataFromDeviceSt struct {
	LightTurnedOn *bool `bson:"lightTurnedOn,omitempty" json:"lightTurnedOn,omitempty"`
	LightLvl *float64 `bson:"lightLvl,omitempty" json:"lightLvl,omitempty"`
}

type PlantainerLightModuleFromDeviceStateSt struct {
	PlantainerLightModuleStateSt `bson:",inline"`
	LightModuleStateDataFromDeviceSt                 `bson:",inline"`
}

type PlantainerShadowStatePieceFromDeviceSt struct {
	PlantainerShadowStatePieceSt `bson:"inline"`
	LightModule *PlantainerLightModuleFromDeviceStateSt `bson:"lightModule" json:"lightModule,omitempty"`
}

type PlantainerShadowRPCMsgFromDeviceArgsStateSt struct {
	Reported *PlantainerShadowStatePieceFromDeviceSt `json:"reported,omitempty"`
	Desired *PlantainerShadowStatePieceFromDeviceSt `json:"desired,omitempty"`
}

type PlantainerShadowRPCMsgFromDeviceArgsSt struct {
	State     PlantainerShadowRPCMsgFromDeviceArgsStateSt
	Version   int
}

//easyjson:json
type ShadowUpdateRPCMsgFromDeviceSt struct {
	mccommunication.RPCMsg
	Args PlantainerShadowRPCMsgFromDeviceArgsSt
}

func (this *ShadowUpdateRPCMsgFromDeviceSt) ConvertToShadowUpdateRPCMsgSt() *ShadowUpdateRPCMsgSt {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("ConvertToShadowUpdateRPCMsgSt recovered: ", r)
		}
	}()

	res := &ShadowUpdateRPCMsgSt{
		RPCMsg: this.RPCMsg,
		Args: PlantainerShadowRPCMsgArgsSt{
			Version: this.Args.Version,
			State: PlantainerShadowRPCMsgArgsStateSt{},
		},
	}
	if this.Args.State.Reported != nil {
		res.Args.State.Reported = &PlantainerShadowStatePieceSt{
			VentilationModule: this.Args.State.Reported.VentilationModule,
			IrrigationModule: this.Args.State.Reported.IrrigationModule,
			// !Add new module functions here!
		}
		if this.Args.State.Reported.LightModule != nil {
			res.Args.State.Reported.LightModule = &this.Args.State.Reported.LightModule.PlantainerLightModuleStateSt
			if this.Args.State.Reported.LightModule.LightLvl != nil {
				in := int(*this.Args.State.Reported.LightModule.LightLvl)
				res.Args.State.Reported.LightModule.LightLvl = &in
			}
		}
	}
	if this.Args.State.Desired != nil {
		res.Args.State.Desired = &PlantainerShadowStatePieceSt{
			VentilationModule: this.Args.State.Desired.VentilationModule,
			IrrigationModule: this.Args.State.Desired.IrrigationModule,
			// !Add new module functions here!
		}
		if this.Args.State.Desired.LightModule != nil {
			res.Args.State.Desired.LightModule = &this.Args.State.Desired.LightModule.PlantainerLightModuleStateSt
			if this.Args.State.Desired.LightModule.LightLvl != nil {
				in := int(*this.Args.State.Desired.LightModule.LightLvl)
				res.Args.State.Desired.LightModule.LightLvl = &in
			}
		}
	}
	return res
}
