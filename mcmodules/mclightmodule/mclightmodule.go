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

//type LightModuleSt struct {
//	State *LightModuleStateSt `bson:"-"`
//	Timer *TimerSt `bson:"-"`
//	DeviceId string `bson:"-"`
//}
//
//func NewLightModule() *LightModuleSt {
//	return &LightModuleSt{
//		NewLightModuleStateWithDefault(),
//		NewTimer(),
//		"",
//	}
//}
//
//func (this *LightModuleSt) SetState(state *LightModuleStateSt) {
//	this.State = state
//}
//
//func (this *LightModuleSt) CheckState() {
//	now := time.Now()
//	nowHour := now.Hour()
//	nowMin := now.Minute()
//	for _, interval := range *this.State.LightIntervalsArr {
//		if nowHour >= interval.FromTimeHours && nowMin >= interval.FromTimeMinutes && nowHour <= interval.ToTimeHours && nowMin < interval.ToTimeMinutes {
//			this.State.LightTurnedOn = interval.TurnedOn
//			// Set Desired state and send Delta
//		}
//	}
//}
//
//func (this *LightModuleSt) SetTimer() {
//	this.Timer.Set(*this.State.LightIntervalsCheckingInterval)
//	<- this.Timer.Timer.C
//	this.Timer.isRunning = false
//	this.CheckState()
//}
//
//func (this *LightModuleSt) ResetTimer() {
//	this.Timer.Reset(*this.State.LightIntervalsCheckingInterval)
//}
//
//func (this *LightModuleSt) RemoveTimer() {
//	this.Timer.Remove()
//}
//
//type TimerSt struct {
//	isRunning bool
//	Timer *time.Timer
//}
//
//func NewTimer() *TimerSt {
//	return &TimerSt{
//		isRunning: false,
//	}
//}
//
//func (this *TimerSt) Set(timeout int) {
//	this.Remove()
//	this.Timer = time.NewTimer(time.Millisecond * time.Duration(timeout))
//	this.isRunning = true
//}
//
//func (this *TimerSt) Reset(timeout int) {
//	if this.Timer != nil && this.isRunning {
//		this.Timer.Reset(time.Millisecond * time.Duration(timeout))
//		this.isRunning = true
//	}
//}
//
//func (this *TimerSt) Remove() {
//	if this.Timer != nil && this.isRunning {
//		this.Timer.Stop()
//	}
//}
//
//func (this *LightModuleSt) CheckOnStateUpdate(deviceId string, newState *LightModuleStateSt) {
//	this.DeviceId = deviceId
//
//	if newState.LightIntervalsCheckingInterval != nil {
//		this.State.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
//		// Reset timers
//		if newState.Mode == nil || newState.Mode == this.State.Mode {
//			//this.ResetTimer()
//		}
//	}
//	if newState.LightIntervalsRestTimeTurnedOn != nil {
//		this.State.LightIntervalsRestTimeTurnedOn = newState.LightIntervalsRestTimeTurnedOn
//	}
//	if newState.LightIntervalsCheckingInterval != nil {
//		this.State.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
//	}
//	if newState.LightLvlCheckActive != nil {
//		this.State.LightLvlCheckActive = newState.LightLvlCheckActive
//	}
//	if newState.Mode != nil && newState.Mode != this.State.Mode {
//		this.State.Mode = newState.Mode
//		switch *newState.Mode {
//		case LightModuleModes[LightModuleModeLightServerIntervalsTimerMode]:
//			// Add timers
//			//this.SetTimer()
//		case LightModuleModes[LightModuleModeManual]:
//			// Reset timers
//			//this.ResetTimer()
//		}
//	}
//}
