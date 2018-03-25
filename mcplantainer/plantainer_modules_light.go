package mcplantainer

import "mevericcore/mclightmodule"

type PlantainerLightModuleStateSt struct {
	mclightmodule.LightModuleStateSt `bson:",inline"`
}

func NewPlantainerLightModuleStateSt() *PlantainerLightModuleStateSt {
	var lightIntervalsArr []mclightmodule.LightModuleInterval = nil
	mode := "manual"
	lightTurnedOn := false
	lightLvlCheckActive := false
	lightLvlCheckInterval := 5100
	lightIntervalsRestTimeTurnedOn := false
	lightIntervalsCheckingInterval := 20000
	return &PlantainerLightModuleStateSt{
		mclightmodule.LightModuleStateSt{
			Mode: &mode,
			LightTurnedOn: &lightTurnedOn,
			LightLvlCheckActive: &lightLvlCheckActive,
			LightLvlCheckInterval: &lightLvlCheckInterval,
			LightIntervalsArr: &lightIntervalsArr,
			LightIntervalsRestTimeTurnedOn: &lightIntervalsRestTimeTurnedOn,
			LightIntervalsCheckingInterval: &lightIntervalsCheckingInterval,
		},
	}
}


func (this *PlantainerLightModuleStateSt) ReportedUpdate(newState *PlantainerLightModuleStateSt) {
	if newState.LightIntervalsCheckingInterval != nil {
		this.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
		// Reset timers
		if newState.Mode == nil || newState.Mode == this.Mode {
			//this.ResetTimer()
		}
	}
	if newState.LightIntervalsRestTimeTurnedOn != nil {
		this.LightIntervalsRestTimeTurnedOn = newState.LightIntervalsRestTimeTurnedOn
	}
	if newState.LightLvlCheckActive != nil {
		this.LightLvlCheckActive = newState.LightLvlCheckActive
	}
	if newState.Mode != nil && newState.Mode != this.Mode {
		this.Mode = newState.Mode
		switch *newState.Mode {
		case mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode]:
			// Add timers
			//this.SetTimer()
		case mclightmodule.LightModuleModes[mclightmodule.LightModuleModeManual]:
			// Reset timers
			//this.ResetTimer()
		}
	}
}

func (this *PlantainerLightModuleStateSt) DesiredUpdate(newState *PlantainerLightModuleStateSt) {
	if newState.LightIntervalsCheckingInterval != nil {
		this.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
	}
	if newState.LightIntervalsRestTimeTurnedOn != nil {
		this.LightIntervalsRestTimeTurnedOn = newState.LightIntervalsRestTimeTurnedOn
	}
	if newState.LightLvlCheckActive != nil {
		this.LightLvlCheckActive = newState.LightLvlCheckActive
	}
	if newState.Mode != nil && newState.Mode != this.Mode {
		this.Mode = newState.Mode
		switch *newState.Mode {
		case mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode]:
			// Add timers
			//this.SetTimer()
		case mclightmodule.LightModuleModes[mclightmodule.LightModuleModeManual]:
			// Reset timers
			//this.ResetTimer()
		}
	}
}