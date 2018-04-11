package mcplantainer

import (
	"mevericcore/mclibs/mccommon"
	"gopkg.in/mgo.v2"
	"mevericcore/mcmodules/mclightmodule"
	"mevericcore/mcmodules/mcventilationmodule"
)

type PlantainerDataValuesIrrigationModuleSt struct {
	Humidity int
	Temperature int
}

//easyjson:json
type PlantainerDataValuesSt struct {
	IrrigationModule *PlantainerDataValuesIrrigationModuleSt
	LightModule *mclightmodule.LightModuleStateDataSt
	VentilationModule *mcventilationmodule.VentilationModuleStateDataSt
}

func NewPlantainerDataValuesSt() *PlantainerDataValuesSt {
	return &PlantainerDataValuesSt{
		&PlantainerDataValuesIrrigationModuleSt{},
		&mclightmodule.LightModuleStateDataSt{},
		&mcventilationmodule.VentilationModuleStateDataSt{},
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