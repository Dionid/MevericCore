package mcplantainer

import (
	"mevericcore/mcmodules/mclightmodule"
	"github.com/robfig/cron"
	"mevericcore/mccommunication"
	"strconv"
	"fmt"
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

func (this *PlantainerLightModuleStateSt) SetCronTasks(deviceId string, cron *cron.Cron) {
	if *this.Mode != mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
		return
	}
	for _, interval := range *this.LightIntervalsArr {
		fromCrString := "0 " + strconv.Itoa(interval.FromTimeMinutes) + " " + strconv.Itoa(interval.FromTimeHours) + " * * *"
		cron.AddFunc(fromCrString, func() {
			// ToDo: change Method to normal way
			fmt.Println("Cron " + strconv.Itoa(interval.FromTimeMinutes))
			rpcMsg := &mccommunication.RPCMsg{
				Src: deviceId,
				Dst: PlantainerServerId,
				Method: deviceId + ".Shadow.Update",
				Args: map[string]interface{}{
					"state": map[string]interface{}{
						"desired": map[string]interface{}{
							"lightModule": PlantainerLightModuleStateSt{
								mclightmodule.LightModuleStateSt{
									LightModuleStateDataSt: mclightmodule.LightModuleStateDataSt{
										LightTurnedOn: &interval.TurnedOn,
									},
								},
							},
						},
					},
				},
			}
			bData, _ := rpcMsg.MarshalJSON()
			req := &mccommunication.ClientToServerRPCReqSt{
				ClientToServerReqSt: mccommunication.ClientToServerReqSt{
					ClientId: deviceId,
					Protocol: "Cron",
					Msg: &bData,
				},
				RPCMsg: rpcMsg,
			}
			// ToDo: CHANGE THIS SHIT FROM ClientToServerRPCReqSt TO RPCMsg
			innerRPCMan.PublishClientToServerRPCReq("Devices.Plantainer.Cron.Task.Exec", req)
		})
		toCrStr := "0 " + strconv.Itoa(interval.ToTimeMinutes) + " " + strconv.Itoa(interval.ToTimeHours) + " * * *"
		cron.AddFunc(toCrStr, func() {
			// ToDo: change Method to normal way
			rpcMsg := &mccommunication.RPCMsg{
				Src: deviceId,
				Dst: PlantainerServerId,
				Method: deviceId + ".Shadow.Update",
				Args: map[string]interface{}{
					"state": map[string]interface{}{
						"desired": map[string]interface{}{
							"lightModule": PlantainerLightModuleStateSt{
								mclightmodule.LightModuleStateSt{
									LightModuleStateDataSt: mclightmodule.LightModuleStateDataSt{
										LightTurnedOn: this.LightIntervalsRestTimeTurnedOn,
									},
								},
							},
						},
					},
				},
			}
			bData, _ := rpcMsg.MarshalJSON()
			req := &mccommunication.ClientToServerRPCReqSt{
				ClientToServerReqSt: mccommunication.ClientToServerReqSt{
					ClientId: deviceId,
					Protocol: "Cron",
					Msg: &bData,
				},
				RPCMsg: rpcMsg,
			}
			innerRPCMan.PublishClientToServerRPCReq("Devices.Plantainer.Cron.Task.Exec", req)
		})
	}
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
	if needToResetTimers || needToStopTimers {
		method := "DeviceCron.Plantainer.RPC.Reset"
		if needToStopTimers {
			method = "DeviceCron.Plantainer.RPC.Stop"
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
		innerRPCMan.PublishRPC("DeviceCron.Plantainer.RPC", rpcMsg)
	}
}

func (this *PlantainerLightModuleStateSt) ReportedUpdate(newState *PlantainerLightModuleStateSt) {
	if newState.LightTurnedOn != nil {
		this.LightTurnedOn = newState.LightTurnedOn
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
		// Todo: add validation
		this.LightIntervalsArr = newState.LightIntervalsArr
	}

	if newState.LightIntervalsCheckingInterval != nil {
		this.LightIntervalsCheckingInterval = newState.LightIntervalsCheckingInterval
	}

	if newState.LightIntervalsRestTimeTurnedOn != nil {
		this.LightIntervalsRestTimeTurnedOn = newState.LightIntervalsRestTimeTurnedOn
	}
}