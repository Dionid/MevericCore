package mcplantainer

import (
	"mevericcore/mcmodules/mcirrigationmodule"
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

// This function must check if it's time to turn on irrigation system
func (this *PlantainerIrrigationModuleStateSt) CheckAllSystems(desiredState *PlantainerIrrigationModuleStateSt) (changed bool, err error) {
	//changed = false
	//err = nil
	//if this.IrrigationTimerLastCallEndTimestamp != nil {
	//	now := time.Now()
	//	nowMHinS := now.Hour() * 3600000 + now.Minute() * 60000
	//	xS := *this.IrrigationTimerEveryXSeconds
	//	tm := time.Unix(int64(*this.IrrigationTimerLastCallEndTimestamp), 0)
	//	tmH := tm.Hour()
	//	tmM := tm.Minute()
	//	tmMHinS := tmM * 60000 + tmH * 3600000
	//	//Here we've got checking that's enough day time has been passed from last activation
	//	//If it's so we can turn IrrigationTurnedOn
	//	if (tmMHinS + xS) < nowMHinS {
	//		if !*this.IrrigationTurnedOn {
	//			t := true
	//			desiredState.IrrigationTurnedOn = &t
	//			//ToDo: deside if it's nassessary to add
	//			changed = true
	//		}
	//	}
	//} else {
	//	if !*this.IrrigationTurnedOn {
	//		t := true
	//		desiredState.IrrigationTurnedOn = &t
	//		changed = true
	//	}
	//}
	return
}

func (this *PlantainerIrrigationModuleStateSt) CheckAfterShadowUpdate(deviceId string, oldState *PlantainerIrrigationModuleStateSt) {
	//needToResetTimers := false
	//needToStopTimers := false
	//
	//if *oldState.Mode != *this.Mode {
	//	switch *this.Mode {
	//	case mcirrigation.LightModuleModes[mcirrigation.LightModuleModeLightServerIntervalsTimerMode]:
	//		// Reset timers
	//		needToResetTimers = true
	//	case mcirrigation.LightModuleModes[mcirrigation.LightModuleModeManual]:
	//		// Stop timers
	//		needToStopTimers = true
	//	}
	//} else {
	//	if oldState.IrrigationTimerEveryXSeconds != this.IrrigationTimerEveryXSeconds ||
	//		oldState.IrrigationTimerIrrigateYSeconds != this.IrrigationTimerIrrigateYSeconds {
	//		if *this.Mode == mcirrigationmodule.IrrigationModuleModeServerIrrigationTimerMode {
	//			// Reset Timer
	//			needToResetTimers = true
	//		}
	//	}
	//}
	//// ToDO: Move this away from here (models must be thin)
	//if needToResetTimers || needToStopTimers {
	//	method := "Plantainer.Cron.Reset"
	//	if needToStopTimers {
	//		method = "Plantainer.Cron.Stop"
	//	}
	//	rpcMsg := &mccommunication.RPCMsg{
	//		Src: deviceId,
	//		Dst: PlantainerServerId,
	//		Method: method,
	//		Args: map[string]interface{}{
	//			"deviceId": deviceId,
	//			"modules": []string{"irrigationModule"},
	//		},
	//	}
	//	innerRPCMan.PublishRPC("Plantainer.Cron.RPC", rpcMsg)
	//}
}