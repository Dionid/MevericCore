package mcplantainer

import (
	"mevericcore/mcmodules/mcirrigationmodule"
	"mevericcore/mcmodules/mclightmodule"
	"time"
	"mevericcore/mclibs/mccommunication"
)

type PlantainerIrrigationModuleStateSt struct {
	mcirrigationmodule.IrrigationModuleStateSt `bson:",inline"`
}

func NewPlantainerIrrigationModuleStateSt() *PlantainerIrrigationModuleStateSt {
	Mode := mcirrigationmodule.IrrigationModuleModeManual
	HumidityCheckActive := false
	HumidityCheckInterval := 20000
	HumidityCheckMinLvl := 50
	HumidityCheckAverageLvl := 60
	HumidityCheckMaxLvl := 100
	IrrigationTimerInProgress := false
	IrrigationTimerEveryXSeconds := 300000
	IrrigationTimerIrrigateYSeconds := 15000
	IrrigationTurnedOn := false

	return &PlantainerIrrigationModuleStateSt{
		mcirrigationmodule.IrrigationModuleStateSt{
			IrrigationModuleStateDataSt: mcirrigationmodule.IrrigationModuleStateDataSt{
				IrrigationTurnedOn: &IrrigationTurnedOn,
			},
			Mode: &Mode,
			HumidityCheckActive: &HumidityCheckActive,
			HumidityCheckInterval: &HumidityCheckInterval,
			HumidityCheckMinLvl: &HumidityCheckMinLvl,
			HumidityCheckAverageLvl: &HumidityCheckAverageLvl,
			HumidityCheckMaxLvl: &HumidityCheckMaxLvl,
			IrrigationTimerInProgress: &IrrigationTimerInProgress,
			IrrigationTimerEveryXSeconds: &IrrigationTimerEveryXSeconds,
			IrrigationTimerIrrigateYSeconds: &IrrigationTimerIrrigateYSeconds,
		},
	}
}

func (this *PlantainerIrrigationModuleStateSt) CheckAllSystems(desiredState *PlantainerIrrigationModuleStateSt) (changed bool, err error) {
	changed = false
	err = nil
	now := time.Now()
	//ToDo: From int to Timestamp and compare with now
	if (*this.IrrigationTimerLastCallEndTimestamp + *this.IrrigationTimerEveryXSeconds) < now {
		if !*this.IrrigationTurnedOn {
			t := true
			desiredState.IrrigationTurnedOn = &t
			changed = true
		}
	}
	return
}

func (this *PlantainerIrrigationModuleStateSt) CheckAfterShadowUpdate(deviceId string, oldState *PlantainerIrrigationModuleStateSt) {
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
		if oldState.IrrigationTimerEveryXSeconds != this.IrrigationTimerEveryXSeconds ||
			oldState.IrrigationTimerIrrigateYSeconds != this.IrrigationTimerIrrigateYSeconds {
			if *this.Mode == mcirrigationmodule.IrrigationModuleModeServerIrrigationTimerMode {
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
				"modules": []string{"irrigationModule"},
			},
		}
		innerRPCMan.PublishRPC("Plantainer.Cron.RPC", rpcMsg)
	}
}