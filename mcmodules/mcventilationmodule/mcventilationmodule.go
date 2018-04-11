package mcventilationmodule

var VentilationModuleModeManual = "manual"
var VentilationModuleModeAuto = "auto"
var VentilationModuleServerCheck = "serverCheck"

var VentilationModuleModes = map[string]string{
	VentilationModuleModeManual:                        "manual",
	VentilationModuleModeAuto: "auto",
	VentilationModuleServerCheck: "serverCheck",
}

type VentilationModuleStateDataSt struct {
	Humidity *float64
	CoolerInTurnedOn *bool `bson:"coolerInTurnedOn,omitempty"`
	CoolerOutTurnedOn *bool `bson:"coolerOutTurnedOn,omitempty"`
}

type VentilationModuleStateSt struct {
	VentilationModuleStateDataSt                 `bson:",inline"`
	Mode                                   *string                `bson:"mode,omitempty"`
	Interval *int
	HumidityMaxLvl *int `bson:"humidityMaxLvl,omitempty"`
	HumidityAverageLvl *int `bson:"humidityAverageLvl,omitempty"`
}

func (this *VentilationModuleStateSt) DesiredUpdate(newState *VentilationModuleStateSt) error {
	if newState.Mode != nil {
		this.Mode = newState.Mode
	}
	if newState.Interval != nil {
		this.Interval = newState.Interval
	}
	if newState.HumidityMaxLvl != nil {
		this.HumidityMaxLvl = newState.HumidityMaxLvl
	}
	if newState.HumidityAverageLvl != nil {
		this.HumidityAverageLvl = newState.HumidityAverageLvl
	}
	if newState.CoolerInTurnedOn != nil {
		this.CoolerInTurnedOn = newState.CoolerInTurnedOn
	}
	if newState.CoolerOutTurnedOn != nil {
		this.CoolerOutTurnedOn = newState.CoolerOutTurnedOn
	}
	return nil
}

func (this *VentilationModuleStateSt) ReportedUpdate(newState *VentilationModuleStateSt) {
	if newState.Mode != nil && newState.Mode != this.Mode {
		this.Mode = newState.Mode
		//switch *this.Mode {
		//case VentilationModuleModes[VentilationModuleModeManual]:
		//	stop = true
		//case VentilationModuleModes[VentilationModuleServerCheck]:
		//	reset = true
		//}
	}
	if newState.Interval != nil {
		this.Interval = newState.Interval
		//if *this.Mode == VentilationModuleModeAuto {
		//	reset = true
		//}
	}
	if newState.Humidity != nil {
		this.Humidity = newState.Humidity
	}
	if newState.CoolerInTurnedOn != nil {
		this.CoolerInTurnedOn = newState.CoolerInTurnedOn
	}
	if newState.CoolerOutTurnedOn != nil {
		this.CoolerOutTurnedOn = newState.CoolerOutTurnedOn
	}
	if newState.HumidityMaxLvl != nil {
		this.HumidityMaxLvl = newState.HumidityMaxLvl
	}
	if newState.HumidityAverageLvl != nil {
		this.HumidityAverageLvl = newState.HumidityAverageLvl
	}
	//return stop, reset
}