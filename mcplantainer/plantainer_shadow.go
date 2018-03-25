package mcplantainer

import (
	"mevericcore/mccommunication"
	"encoding/json"
	"fmt"
)

//easyjson:json
type PlantainerShadowStatePieceSt struct {
	LightModule PlantainerLightModuleStateSt `bson:"lightModule"`
}

func NewPlantainerShadowStatePiece() *PlantainerShadowStatePieceSt {
	return &PlantainerShadowStatePieceSt{
		*NewPlantainerLightModuleStateSt(),
	}
}

//easyjson:json
type PlantainerShadowStateSt struct {
	Reported PlantainerShadowStatePieceSt
	Desired *PlantainerShadowStatePieceSt
	Delta *PlantainerShadowStatePieceSt `bson:"-"`
}

func NewPlantainerShadowState() *PlantainerShadowStateSt {
	return &PlantainerShadowStateSt{
		*NewPlantainerShadowStatePiece(),
		nil,
		nil,
	}
}

func (this *PlantainerShadowStateSt) fillDelta(reported *map[string]interface{}, desired *map[string]interface{}, delta *map[string]interface{}) {
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
				if (*reported)[key] != val && val != nil {
					(*delta)[key] = val
				}
			}
		}
	}
}

func (this *PlantainerShadowStateSt) FillDelta() *map[string]interface{} {
	des := this.Desired
	if des == nil {
		return nil
	}

	//_______

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

type PlantainerShadowMetadataSt struct {
	Version int
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
