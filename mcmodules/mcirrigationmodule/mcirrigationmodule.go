package mcirrigationmodule

var IrrigationModuleModeManual = "manual"
var IrrigationModuleModeDeviceIrrigationTimerMode = "deviceIrrigationTimerMode"
var IrrigationModuleModeServerIrrigationTimerMode = "serverIrrigationTimerMode"
var IrrigationModuleModeHumidityMode = "humidityMode"

type IrrigationModuleStateDataSt struct {
	IrrigationTurnedOn *bool `bson:"irrigationTurnedOn,omitempty"`
	Humidity *float64 `bson:"humidity,omitempty"`
}

type IrrigationModuleStateSt struct {
	IrrigationModuleStateDataSt  `bson:",inline"`

	Mode *string `bson:"mode,omitempty"`

	// humidityCheck
	HumidityCheckActive *bool `bson:"humidityCheckActive,omitempty"`
	HumidityCheckInterval *int `bson:"humidityCheckInterval,omitempty"`
	HumidityCheckLastIntervalCallTimestamp *int `bson:"humidityCheckLastIntervalCallTimestamp,omitempty"`
	HumidityCheckMinLvl *int `bson:"humidityCheckMinLvl,omitempty"`
	HumidityCheckAverageLvl *int `bson:"humidityCheckAverageLvl,omitempty"`
	HumidityCheckMaxLvl *int `bson:"humidityCheckMaxLvl,omitempty"`

	// irrigationTimer
	IrrigationTimerInProgress *bool `bson:"irrigationTimerInProgress,omitempty"`
	IrrigationTimerEveryXSeconds *int `bson:"irrigationTimerEveryXSeconds,omitempty"`
	IrrigationTimerIrrigateYSeconds *int `bson:"irrigationTimerIrrigateYSeconds,omitempty"`
	IrrigationTimerLastCallStartTimestamp *int `bson:"irrigationTimerLastCallStartTimestamp,omitempty"`
	IrrigationTimerLastCallEndTimestamp *int `bson:"irrigationTimerLastCallEndTimestamp,omitempty"`
}

func (this *IrrigationModuleStateSt) DesiredUpdate(newState *IrrigationModuleStateSt) error {
	if newState.Mode != nil {
		this.Mode = newState.Mode
	}
	if newState.HumidityCheckActive != nil {
		this.HumidityCheckActive = newState.HumidityCheckActive
	}
	if newState.HumidityCheckInterval != nil {
		this.HumidityCheckInterval = newState.HumidityCheckInterval
	}
	if newState.HumidityCheckMinLvl != nil {
		this.HumidityCheckMinLvl = newState.HumidityCheckMinLvl
	}
	if newState.HumidityCheckAverageLvl != nil {
		this.HumidityCheckAverageLvl = newState.HumidityCheckAverageLvl
	}
	if newState.HumidityCheckMaxLvl != nil {
		this.HumidityCheckMaxLvl = newState.HumidityCheckMaxLvl
	}
	if newState.IrrigationTurnedOn != nil {
		this.IrrigationTurnedOn = newState.IrrigationTurnedOn
	}
	if newState.IrrigationTimerEveryXSeconds != nil {
		this.IrrigationTimerEveryXSeconds = newState.IrrigationTimerEveryXSeconds
	}
	if newState.IrrigationTimerIrrigateYSeconds != nil {
		this.IrrigationTimerIrrigateYSeconds = newState.IrrigationTimerIrrigateYSeconds
	}
	return nil
}


func (this *IrrigationModuleStateSt) ReportedUpdate(newState *IrrigationModuleStateSt) (stop bool, reset bool) {
	if newState.Mode != nil {
		this.Mode = newState.Mode
	}
	if newState.HumidityCheckActive != nil {
		this.HumidityCheckActive = newState.HumidityCheckActive
	}
	if newState.HumidityCheckInterval != nil {
		this.HumidityCheckInterval = newState.HumidityCheckInterval
	}
	if newState.HumidityCheckMinLvl != nil {
		this.HumidityCheckMinLvl = newState.HumidityCheckMinLvl
	}
	if newState.HumidityCheckAverageLvl != nil {
		this.HumidityCheckAverageLvl = newState.HumidityCheckAverageLvl
	}
	if newState.HumidityCheckMaxLvl != nil {
		this.HumidityCheckMaxLvl = newState.HumidityCheckMaxLvl
	}
	if newState.IrrigationTurnedOn != nil {
		this.IrrigationTurnedOn = newState.IrrigationTurnedOn
	}
	if newState.IrrigationTimerEveryXSeconds != nil {
		this.IrrigationTimerEveryXSeconds = newState.IrrigationTimerEveryXSeconds
	}
	if newState.IrrigationTimerIrrigateYSeconds != nil {
		this.IrrigationTimerIrrigateYSeconds = newState.IrrigationTimerIrrigateYSeconds
	}
	if newState.HumidityCheckLastIntervalCallTimestamp != nil {
		this.HumidityCheckLastIntervalCallTimestamp = newState.HumidityCheckLastIntervalCallTimestamp
	}
	if newState.IrrigationTimerLastCallStartTimestamp != nil {
		this.IrrigationTimerLastCallStartTimestamp = newState.IrrigationTimerLastCallStartTimestamp
	}
	if newState.IrrigationTimerLastCallEndTimestamp != nil {
		this.IrrigationTimerLastCallEndTimestamp = newState.IrrigationTimerLastCallEndTimestamp
	}
	if newState.Mode != nil && newState.Mode != this.Mode {
		this.Mode = newState.Mode
		switch *newState.Mode {
		case IrrigationModuleModeServerIrrigationTimerMode:
			// Reset timers
			return false, true
		case IrrigationModuleModeHumidityMode:
		case IrrigationModuleModeManual:
			// Stop timers
			return true, false
		}
	}
	return false, false
}