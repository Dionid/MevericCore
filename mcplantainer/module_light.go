package mcplantainer

import "time"

var lightModuleModeManual = "Manual"
var lightModuleModeLightServerIntervalsTimerMode = "LightServerIntervalsTimerMode"

var lightModuleModes = map[string]string{
	lightModuleModeManual: "manual",
	lightModuleModeLightServerIntervalsTimerMode: "lightServerIntervalsTimerMode",
}

type LightModuleInterval struct {
	FromTimeHours int
	FromTimeMinutes int
	ToTimeHours int
	ToTimeMinutes int
	TurnedOn bool
}

type LightModuleStateSt struct {
	Mode *string
	LightTurnedOn bool
	LightLvl *int
	LightLvlCheckActive *bool
	LightLvlCheckInterval *int
	LightLvlCheckLastIntervalCallTimestamp *int
	LightIntervalsArr *[]LightModuleInterval
	LightIntervalsRestTimeTurnedOn *bool
	LightIntervalsCheckingInterval *int
}

func NewLightModuleState(mode *string, lightLvlCheckActive *bool, lightLvlCheckInterval *int, lightIntervalsRestTimeTurnedOn *bool, lightIntervalsCheckingInterval *int, lightIntervalsArr *[]LightModuleInterval) *LightModuleStateSt {
	return &LightModuleStateSt{
		Mode: mode,
		LightLvlCheckActive: lightLvlCheckActive,
		LightLvlCheckInterval: lightLvlCheckInterval,
		LightIntervalsArr: lightIntervalsArr,
		LightIntervalsRestTimeTurnedOn: lightIntervalsRestTimeTurnedOn,
		LightIntervalsCheckingInterval: lightIntervalsCheckingInterval,
	}
}

func NewLightModuleStateWithDefault() *LightModuleStateSt {
	var lightIntervalsArr []LightModuleInterval = nil
	mode := "manual"
	lightLvlCheckActive := false
	lightLvlCheckInterval := 5100
	lightIntervalsRestTimeTurnedOn := false
	lightIntervalsCheckingInterval := 20000
	return &LightModuleStateSt{
		Mode: &mode,
		LightLvlCheckActive: &lightLvlCheckActive,
		LightLvlCheckInterval: &lightLvlCheckInterval,
		LightIntervalsArr: &lightIntervalsArr,
		LightIntervalsRestTimeTurnedOn: &lightIntervalsRestTimeTurnedOn,
		LightIntervalsCheckingInterval: &lightIntervalsCheckingInterval,
	}
}

type LightModuleSt struct {
	State *LightModuleStateSt `bson:"-"`
	Timer *TimerSt `bson:"-"`
	DeviceId string `bson:"-"`
}

func NewLightModule() *LightModuleSt {
	return &LightModuleSt{
		NewLightModuleStateWithDefault(),
		NewTimer(),
		"",
	}
}

func (this *LightModuleSt) SetState(state *LightModuleStateSt) {
	this.State = state
}

func (this *LightModuleSt) CheckState() {
	now := time.Now()
	nowHour := now.Hour()
	nowMin := now.Minute()
	for _, interval := range *this.State.LightIntervalsArr {
		if nowHour >= interval.FromTimeHours && nowMin >= interval.FromTimeMinutes && nowHour <= interval.ToTimeHours && nowMin < interval.ToTimeMinutes {
			this.State.LightTurnedOn = interval.TurnedOn
			// Set Desired state and send Delta
		}
	}
}

func (this *LightModuleSt) SetTimer() {
	this.Timer.Set(*this.State.LightIntervalsCheckingInterval)
	<- this.Timer.Timer.C
	this.Timer.isRunning = false
	this.CheckState()
}

func (this *LightModuleSt) ResetTimer() {
	this.Timer.Reset(*this.State.LightIntervalsCheckingInterval)
}

func (this *LightModuleSt) RemoveTimer() {
	this.Timer.Remove()
}

type TimerSt struct {
	isRunning bool
	Timer *time.Timer
}

func NewTimer() *TimerSt {
	return &TimerSt{
		isRunning: false,
	}
}

func (this *TimerSt) Set(timeout int) {
	this.Remove()
	this.Timer = time.NewTimer(time.Millisecond * time.Duration(timeout))
	this.isRunning = true
}

func (this *TimerSt) Reset(timeout int) {
	if this.Timer != nil && this.isRunning {
		this.Timer.Reset(time.Millisecond * time.Duration(timeout))
		this.isRunning = true
	}
}

func (this *TimerSt) Remove() {
	if this.Timer != nil && this.isRunning {
		this.Timer.Stop()
	}
}

func (this *LightModuleSt) CheckOnStateUpdate(deviceId string, newState *LightModuleStateSt) {
	this.DeviceId = deviceId

	if newState.LightIntervalsCheckingInterval != nil {
		this.State.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
		// Reset timers
		if newState.Mode == nil || newState.Mode == this.State.Mode {
			//this.ResetTimer()
		}
	}
	if newState.LightIntervalsRestTimeTurnedOn != nil {
		this.State.LightIntervalsRestTimeTurnedOn = newState.LightIntervalsRestTimeTurnedOn
	}
	if newState.LightIntervalsCheckingInterval != nil {
		this.State.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
	}
	if newState.LightLvlCheckActive != nil {
		this.State.LightLvlCheckActive = newState.LightLvlCheckActive
	}
	if newState.Mode != nil && newState.Mode != this.State.Mode {
		this.State.Mode = newState.Mode
		switch *newState.Mode {
		case lightModuleModes[lightModuleModeLightServerIntervalsTimerMode]:
			// Add timers
			//this.SetTimer()
		case lightModuleModes[lightModuleModeManual]:
			// Reset timers
			//this.ResetTimer()
		}
	}
}