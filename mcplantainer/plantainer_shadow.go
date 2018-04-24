package mcplantainer

import (
	"mevericcore/mclibs/mccommunication"
	"encoding/json"
	"fmt"
	"reflect"
	"mevericcore/mcmodules/mclightmodule"
)

//easyjson:json
type PlantainerShadowStatePieceSt struct {
	LightModule *PlantainerLightModuleStateSt `bson:"lightModule" json:"lightModule,omitempty"`
	VentilationModule *PlantainerVentilationModuleStateSt `bson:"ventilationModule" json:"ventilationModule,omitempty"`
	IrrigationModule *PlantainerIrrigationModuleStateSt `bson:"irrigationModule" json:"irrigationModule,omitempty"`
}

func NewPlantainerShadowStatePiece() *PlantainerShadowStatePieceSt {
	return &PlantainerShadowStatePieceSt{
		NewPlLightModuleStateWithDefaultsSt(),
		NewPlantainerVentilationModuleState(),
		NewPlantainerIrrigationModuleStateSt(),
	}
}

//easyjson:json
type PlantainerShadowStateSt struct {
	Reported PlantainerShadowStatePieceSt
	Desired PlantainerShadowStatePieceSt
	Delta *PlantainerShadowStatePieceSt `bson:"-"`
}

func NewPlantainerShadowState() *PlantainerShadowStateSt {
	return &PlantainerShadowStateSt{
		PlantainerShadowStatePieceSt{},
		*NewPlantainerShadowStatePiece(),
		nil,
	}
}

func (this *PlantainerShadowStateSt) fillDelta(reported *map[string]interface{}, desired *map[string]interface{}, delta *map[string]interface{}) {
	defer func(){
		if err := recover(); err != nil {
			fmt.Println("Recoverd in fillDelta")
		}
	}()

	for key, val := range *desired {
		if (*reported)[key] != nil {
			switch desireV := val.(type) {
			case map[string]interface{}:
				switch repV := (*reported)[key].(type) {
				case map[string]interface{}:
					newMap := map[string]interface{}{}
					this.fillDelta(&repV, &desireV, &newMap)
					if len(newMap) > 0 {
						(*delta)[key] = newMap
					}
				default:
					(*delta)[key] = val
				}
			default:
				switch repV := (*reported)[key].(type) {
				case []interface{}:
					if !reflect.DeepEqual(repV, val) {
						if val != nil {
							(*delta)[key] = val
						}
					}
				default:
					if repV != val {
						if val != nil {
							(*delta)[key] = val
						}
					}
				}
				//if (*reported)[key] != val {
				//	if val != nil {
				//		(*delta)[key] = val
				//	}
				//}
			}
		}
	}
}

func (this *PlantainerShadowStateSt) FillDelta() *map[string]interface{} {
	des := this.Desired

	bData, err := des.MarshalJSON()
	desMap := map[string]interface{}{}
	if err != nil {
		fmt.Printf("bData err: %+v\n", err.Error())
		return nil
	} else {
		if err := json.Unmarshal(bData, &desMap); err != nil {
			fmt.Printf("bData err: %+v\n", err.Error())
			return nil
		}
		//fmt.Printf("desMap: %+v\n", desMap)
	}

	bResData, err := this.Reported.MarshalJSON()
	repMap := map[string]interface{}{}
	if err != nil {
		fmt.Printf("bResData err: %+v\n", err.Error())
		return nil
	} else {
		if err := json.Unmarshal(bResData, &repMap); err != nil {
			fmt.Printf("bResData err: %+v\n", err.Error())
			return nil
		}
		//fmt.Printf("repMap: %+v\n", repMap)
	}
	deltaMap := map[string]interface{}{}
	//fmt.Printf("deltaMap : %+v\n", deltaMap)
	this.fillDelta(&repMap, &desMap, &deltaMap)
	//fmt.Printf("deltaMap : %+v\n", deltaMap)
	//fmt.Printf("Delta: %+v\n", this.Delta)

	if len(deltaMap) == 0 {
		return nil
	}

	if dBData, err := json.Marshal(deltaMap); err != nil {
		fmt.Printf("dBData err: %+v\n", err.Error())
		return nil
	} else {
		this.Delta = &PlantainerShadowStatePieceSt{}
		//fmt.Printf("dBData: %+v\n", dBData)
		if err := json.Unmarshal(dBData, &this.Delta); err != nil {
			//fmt.Printf("bData err: %+v\n", err.Error())
			return nil
		}
		//fmt.Printf("Success\n")
		//fmt.Printf("Delta: %+v\n", this.Delta)
	}

	return nil
}

// ToDo: Add other data (timestamps)
//easyjson:json
type PlantainerShadowMetadataSt struct {
	Version int `json:"version,omitempty"`
}

//easyjson:json
type PlantainerShadowSt struct {
	Id string
	State PlantainerShadowStateSt
	Metadata PlantainerShadowMetadataSt
}

func NewPlantainerShadow(shadowId string) *PlantainerShadowSt {
	return &PlantainerShadowSt{
		shadowId,
		*NewPlantainerShadowState(),
		PlantainerShadowMetadataSt{},
	}
}

func (this *PlantainerShadowSt) CheckVersion(version int) bool {
	return this.Metadata.Version == version
}

func (this *PlantainerShadowSt) IncrementVersion() {
	this.Metadata.Version += 1
}

type PlantainerShadowRPCMsgArgsStateSt struct {
	Reported *PlantainerShadowStatePieceSt `json:"reported"`
	Desired *PlantainerShadowStatePieceSt `json:"desired"`
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

type LightModuleStateDataFromDeviceSt struct {
	LightTurnedOn *bool `bson:"lightTurnedOn,omitempty" json:"lightTurnedOn,omitempty"`
	LightLvl *float64 `bson:"lightLvl,omitempty" json:"lightLvl,omitempty"`
}

type PlantainerLightModuleFromDeviceStateSt struct {
	mclightmodule.LightModuleStateSt `bson:",inline"`
	LightModuleStateDataFromDeviceSt                 `bson:",inline"`
}

type PlantainerShadowStatePieceFromDeviceSt struct {
	LightModule PlantainerLightModuleFromDeviceStateSt `bson:"lightModule" json:"lightModule,omitempty"`
	VentilationModule PlantainerVentilationModuleStateSt `bson:"ventilationModule" json:"ventilationModule,omitempty"`
	IrrigationModule PlantainerIrrigationModuleStateSt`bson:"irrigationModule" json:"irrigationModule,omitempty"`
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
type JSONShadowUpdateRPCMsgFromDeviceSt struct {
	mccommunication.RPCMsg
	Args PlantainerShadowRPCMsgFromDeviceArgsSt
}

func (this *JSONShadowUpdateRPCMsgFromDeviceSt) ConvertToShadowUpdateRPCMsgSt() *ShadowUpdateRPCMsgSt {
	defer func() {
		r := recover()
		if r != nil {
			print(r)
		}
	}()
	res := &ShadowUpdateRPCMsgSt{
		RPCMsg: this.RPCMsg,
		Args: PlantainerShadowRPCMsgArgsSt{
			Version: this.Args.Version,
			State: PlantainerShadowRPCMsgArgsStateSt{},
		},
	}
	// ToDo: This is fucking bullshit
	if this.Args.State.Reported != nil {
		res.Args.State.Reported = &PlantainerShadowStatePieceSt{
			VentilationModule: &this.Args.State.Reported.VentilationModule,
			IrrigationModule: &this.Args.State.Reported.IrrigationModule,
			LightModule: &PlantainerLightModuleStateSt{
				LightModuleStateSt: mclightmodule.LightModuleStateSt{
					LightModuleStateDataSt: mclightmodule.LightModuleStateDataSt{
						LightTurnedOn: this.Args.State.Reported.LightModule.LightTurnedOn,
					},
					Mode: this.Args.State.Reported.LightModule.Mode,
					LightLvlCheckActive: this.Args.State.Reported.LightModule.LightLvlCheckActive,
					LightLvlCheckInterval: this.Args.State.Reported.LightModule.LightLvlCheckInterval,
					LightLvlCheckLastIntervalCallTimestamp: this.Args.State.Reported.LightModule.LightLvlCheckLastIntervalCallTimestamp,
					LightIntervalsArr: this.Args.State.Reported.LightModule.LightIntervalsArr,
					LightIntervalsRestTimeTurnedOn: this.Args.State.Reported.LightModule.LightIntervalsRestTimeTurnedOn,
					LightIntervalsCheckingInterval: this.Args.State.Reported.LightModule.LightIntervalsCheckingInterval,
				},
			},
		}
		if this.Args.State.Reported.LightModule.LightLvl != nil {
			in := int(*this.Args.State.Reported.LightModule.LightLvl)
			res.Args.State.Reported.LightModule.LightLvl = &in
		}
	}
	if this.Args.State.Desired != nil {
		res.Args.State.Desired = &PlantainerShadowStatePieceSt{
			VentilationModule: &this.Args.State.Desired.VentilationModule,
			IrrigationModule: &this.Args.State.Desired.IrrigationModule,
			LightModule: &PlantainerLightModuleStateSt{
				LightModuleStateSt: mclightmodule.LightModuleStateSt{
					LightModuleStateDataSt: mclightmodule.LightModuleStateDataSt{
						LightTurnedOn: this.Args.State.Desired.LightModule.LightTurnedOn,
					},
					Mode: this.Args.State.Desired.LightModule.Mode,
					LightLvlCheckActive: this.Args.State.Desired.LightModule.LightLvlCheckActive,
					LightLvlCheckInterval: this.Args.State.Desired.LightModule.LightLvlCheckInterval,
					LightLvlCheckLastIntervalCallTimestamp: this.Args.State.Desired.LightModule.LightLvlCheckLastIntervalCallTimestamp,
					LightIntervalsArr: this.Args.State.Desired.LightModule.LightIntervalsArr,
					LightIntervalsRestTimeTurnedOn: this.Args.State.Desired.LightModule.LightIntervalsRestTimeTurnedOn,
					LightIntervalsCheckingInterval: this.Args.State.Desired.LightModule.LightIntervalsCheckingInterval,
				},
			},
		}
		if this.Args.State.Desired.LightModule.LightLvl != nil {
			in := int(*this.Args.State.Desired.LightModule.LightLvl)
			res.Args.State.Desired.LightModule.LightLvl = &in
		}
	}
	return res
}