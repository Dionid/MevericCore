package mclightmodule

//import "time"

var LightModuleModeManual = "manual"
var LightModuleModeLightServerIntervalsTimerMode = "lightServerIntervalsTimerMode"

var LightModuleModes = map[string]string{
	LightModuleModeManual:                        "manual",
	LightModuleModeLightServerIntervalsTimerMode: "lightServerIntervalsTimerMode",
}

type LightModuleInterval struct {
	FromTimeHours   int  `bson:"fromTimeHours,omitempty"`
	FromTimeMinutes int  `bson:"fromTimeMinutes,omitempty"`
	ToTimeHours     int  `bson:"toTimeHours,omitempty"`
	ToTimeMinutes   int  `bson:"toTimeMinutes,omitempty"`
	TurnedOn        bool `bson:"turnedOn,omitempty"`
}

type LightModuleStateDataSt struct {
	LightTurnedOn *bool `bson:"lightTurnedOn,omitempty" json:"lightTurnedOn,omitempty"`
	LightLvl *int `bson:"lightLvl,omitempty" json:"lightLvl,omitempty"`
}

type LightModuleStateSt struct {
	LightModuleStateDataSt                 `bson:",inline"`
	Mode                                   *string                `bson:"mode,omitempty"`
	LightLvlCheckActive                    *bool                  `bson:"lightLvlCheckActive,omitempty"`
	LightLvlCheckInterval                  *int                   `bson:"lightLvlCheckInterval,omitempty"`
	LightLvlCheckLastIntervalCallTimestamp *int                   `bson:"lightLvlCheckLastIntervalCallTimestamp,omitempty"`
	LightIntervalsArr                      *[]LightModuleInterval `bson:"lightIntervalsArr,omitempty"`
	LightIntervalsRestTimeTurnedOn         *bool                  `bson:"lightIntervalsRestTimeTurnedOn,omitempty"`
	LightIntervalsCheckingInterval         *int                   `bson:"lightIntervalsCheckingInterval,omitempty"`
}

func NewLightModuleState(mode string, lightTurnedOn bool, lightLvlCheckActive bool, lightLvlCheckInterval int, lightIntervalsRestTimeTurnedOn bool, lightIntervalsCheckingInterval int, lightIntervalsArr []LightModuleInterval) *LightModuleStateSt {
	return &LightModuleStateSt{
		LightModuleStateDataSt: LightModuleStateDataSt{
			LightTurnedOn:                  &lightTurnedOn,
		},
		Mode:                           &mode,
		LightLvlCheckActive:            &lightLvlCheckActive,
		LightLvlCheckInterval:          &lightLvlCheckInterval,
		LightIntervalsArr:              &lightIntervalsArr,
		LightIntervalsRestTimeTurnedOn: &lightIntervalsRestTimeTurnedOn,
		LightIntervalsCheckingInterval: &lightIntervalsCheckingInterval,
	}
}

func NewLightModuleStateWithDefault() *LightModuleStateSt {
	var lightIntervalsArr []LightModuleInterval = nil
	mode := "manual"
	lightTurnedOn := false
	lightLvlCheckActive := false
	lightLvlCheckInterval := 5100
	lightIntervalsRestTimeTurnedOn := false
	lightIntervalsCheckingInterval := 20000
	return &LightModuleStateSt{
		LightModuleStateDataSt: LightModuleStateDataSt{
			LightTurnedOn:                  &lightTurnedOn,
		},
		Mode:                           &mode,
		LightLvlCheckActive:            &lightLvlCheckActive,
		LightLvlCheckInterval:          &lightLvlCheckInterval,
		LightIntervalsArr:              &lightIntervalsArr,
		LightIntervalsRestTimeTurnedOn: &lightIntervalsRestTimeTurnedOn,
		LightIntervalsCheckingInterval: &lightIntervalsCheckingInterval,
	}
}

func (this *LightModuleStateSt) DesiredUpdate(newState *LightModuleStateSt) error {
	if newState.Mode != nil {
		this.Mode = newState.Mode
	}
	if newState.LightTurnedOn != nil {
		this.LightTurnedOn = newState.LightTurnedOn
	}
	if newState.LightLvlCheckActive != nil {
		this.LightLvlCheckActive = newState.LightLvlCheckActive
	}
	if newState.LightLvlCheckInterval != nil {
		this.LightLvlCheckInterval = newState.LightLvlCheckInterval
	}
	if newState.LightIntervalsArr != nil {
		// ToDo: Add validation of intervals quality
		this.LightIntervalsArr = newState.LightIntervalsArr
	}
	if newState.LightIntervalsRestTimeTurnedOn != nil {
		this.LightIntervalsRestTimeTurnedOn = newState.LightIntervalsRestTimeTurnedOn
	}
	if newState.LightIntervalsCheckingInterval != nil {
		this.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
	}
	return nil
}


// ToDo: Change this to return stop bool, reset bool
func (this *LightModuleStateSt) ReportedUpdate(newState *LightModuleStateSt) {
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
	if newState.LightIntervalsCheckingInterval != nil {
		this.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
	}
	if newState.LightLvlCheckActive != nil {
		this.LightLvlCheckActive = newState.LightLvlCheckActive
	}
	if newState.Mode != nil && newState.Mode != this.Mode {
		this.Mode = newState.Mode
		switch *newState.Mode {
		case LightModuleModes[LightModuleModeLightServerIntervalsTimerMode]:
			// Add timers
			//this.SetTimer()
		case LightModuleModes[LightModuleModeManual]:
			// Reset timers
			//this.ResetTimer()
		}
	}
}