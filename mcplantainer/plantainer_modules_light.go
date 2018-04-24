package mcplantainer

import (
	"mevericcore/mcmodules/mclightmodule"
	"mevericcore/mclibs/mccommunication"
	"time"
)

type PlantainerLightModuleStateSt struct {
	mclightmodule.LightModuleStateSt `bson:",inline"`
}

func NewPlLightModuleStateWithDefaultsSt() *PlantainerLightModuleStateSt {
	var lightIntervalsArr []mclightmodule.LightModuleInterval = nil
	mode := "manual"
	lightTurnedOn := false
	lightLvlCheckActive := false
	lightLvlCheckInterval := 5100
	lightIntervalsRestTimeTurnedOn := false
	lightIntervalsCheckingInterval := 20000
	return &PlantainerLightModuleStateSt{
		mclightmodule.LightModuleStateSt{
			LightModuleStateDataSt: mclightmodule.LightModuleStateDataSt{
				LightTurnedOn:                  &lightTurnedOn,
			},
			Mode: &mode,
			LightLvlCheckActive: &lightLvlCheckActive,
			LightLvlCheckInterval: &lightLvlCheckInterval,
			LightIntervalsArr: &lightIntervalsArr,
			LightIntervalsRestTimeTurnedOn: &lightIntervalsRestTimeTurnedOn,
			LightIntervalsCheckingInterval: &lightIntervalsCheckingInterval,
		},
	}
}

func (this *PlantainerLightModuleStateSt) CheckAllSystems(desiredState *PlantainerLightModuleStateSt) (changed bool, err error) {
	changed = false
	err = nil
	if this.Mode != nil && *this.Mode == mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
		now := time.Now()
		nowH := now.Hour()
		nowM := now.Minute()
		for _, interval := range *this.LightIntervalsArr {
			if interval.FromTimeHours < nowH &&
				interval.FromTimeMinutes < nowM &&
				interval.ToTimeHours > nowH &&
				interval.ToTimeMinutes > nowM {
					if *this.LightTurnedOn != interval.TurnedOn {
						desiredState.LightTurnedOn = &interval.TurnedOn
						changed = true
					}
				}
		}
	}
	return
}

func (this *PlantainerLightModuleStateSt) CheckAfterShadowUpdate(deviceId string, oldState *PlantainerLightModuleStateSt) {
	needToResetTimers := false
	needToStopTimers := false
	if *oldState.Mode != *this.Mode {
		switch *this.Mode {
		case mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode]:
			// Reset timers
			needToResetTimers = true
		case mclightmodule.LightModuleModes[mclightmodule.LightModuleModeManual]:
			// Stop timers
			needToStopTimers = true
		}
	} else {
		if oldState.LightIntervalsArr != this.LightIntervalsArr {
			if *this.Mode == mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
				// Reset Timer
				needToResetTimers = true
			}
		}
		if oldState.LightIntervalsCheckingInterval != this.LightIntervalsCheckingInterval {
			if *this.Mode == mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
				// Reset Timer
				needToResetTimers = true
			}
		}
	}
	// ToDO: Move this away from here (models must be thin)
	if needToResetTimers || needToStopTimers {
		method := "Plantainer.Cron.Reset"
		if needToStopTimers {
			method = "Plantainer.Cron.Stop"
		}
		rpcMsg := &mccommunication.RPCMsg{
			Src: deviceId,
			Dst: PlantainerServerId,
			Method: method,
			Args: map[string]interface{}{
				"deviceId": deviceId,
				"modules": []string{"lightModule"},
			},
		}
		innerRPCMan.PublishRPC("Plantainer.Cron.RPC", rpcMsg)
	}
}

func (this *PlantainerLightModuleStateSt) ReportedUpdate(newState *PlantainerLightModuleStateSt) {
	if newState.LightTurnedOn != nil {
		this.LightTurnedOn = newState.LightTurnedOn
	}

	if newState.LightLvl != nil {
		this.LightLvl = newState.LightLvl
	}

	// Device lvl
	if newState.LightLvlCheckActive != nil {
		this.LightLvlCheckActive = newState.LightLvlCheckActive
	}
	if newState.LightLvlCheckInterval != nil {
		this.LightLvlCheckInterval = newState.LightLvlCheckInterval
	}

	if newState.Mode != nil {
		this.Mode = newState.Mode
	}

	if newState.LightIntervalsArr != nil {
		this.LightIntervalsArr = newState.LightIntervalsArr
	}

	if newState.LightIntervalsCheckingInterval != nil {
		this.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
	}

	if newState.LightIntervalsRestTimeTurnedOn != nil {
		this.LightIntervalsRestTimeTurnedOn = newState.LightIntervalsRestTimeTurnedOn
	}
}