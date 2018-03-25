package mcplantainer

import (
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2"
	"mevericcore/mclightmodule"
)

type PlantainerDataValuesIrrigationModuleSt struct {
	Humidity int
	Temperature int
}

type PlantainerDataValuesSt struct {
	IrrigationModule *PlantainerDataValuesIrrigationModuleSt
	LightModule *mclightmodule.LightModuleStateDataSt
}

func NewPlantainerDataValuesSt() *PlantainerDataValuesSt {
	return &PlantainerDataValuesSt{
		&PlantainerDataValuesIrrigationModuleSt{},
		&mclightmodule.LightModuleStateDataSt{},
	}
}

//easyjson:json
type PlantainerDataSt struct {
	mccommon.DeviceDataBaseSt `bson:",inline"`
	Values               map[string]map[string]PlantainerDataValuesSt
}

func NewPlantainerData() *PlantainerDataSt {
	return &PlantainerDataSt{
		mccommon.DeviceDataBaseSt{
			PeriodInSec: 10000,
		},
		nil,
	}
}

func (this *PlantainerDataSt) EnsureIndex(collection *mgo.Collection) error {
	return nil
}