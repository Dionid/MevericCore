package mcplantainer

import (
	"mevericcore/mcmodules/mcventilationmodule"
)

type PlantainerVentilationModuleStateSt struct {
	mcventilationmodule.VentilationModuleStateSt `bson:",inline"`
}

func NewPlantainerVentilationModuleState() *PlantainerVentilationModuleStateSt {
	mode := "manual"
	CoolerInTurnedOn := false
	CoolerOutTurnedOn := false
	Interval := 2000
	HumidityMaxLvl := 27
	HumidityAverageLvl := 23
	var Humidity float64 = 0
	return &PlantainerVentilationModuleStateSt{
		mcventilationmodule.VentilationModuleStateSt{
			VentilationModuleStateDataSt: mcventilationmodule.VentilationModuleStateDataSt{
				Humidity: &Humidity,
				CoolerInTurnedOn: &CoolerInTurnedOn,
				CoolerOutTurnedOn: &CoolerOutTurnedOn,
			},
			Mode: &mode,
			Interval: &Interval,
			HumidityMaxLvl: &HumidityMaxLvl,
			HumidityAverageLvl: &HumidityAverageLvl,
		},
	}
}