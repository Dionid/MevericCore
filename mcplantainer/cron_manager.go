package mcplantainer

import (
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"time"
	"mevericcore/mcmodules/mclightmodule"
	"strconv"
	"mevericcore/mclibs/mccommunication"
)

// This manager is used in Hub to check operations that needs to be done
type DeviceCronManagerSt struct {
	CronByDeviceId map[string]*DeviceCronSt
}

func NewDeviceCronManager() *DeviceCronManagerSt {
	return &DeviceCronManagerSt{
		map[string]*DeviceCronSt{},
	}
}

func (cr *DeviceCronManagerSt) AddDeviceCron(devId string, c *DeviceCronSt) {
	cr.CronByDeviceId[devId] = c
}

func (cr *DeviceCronManagerSt) StartAllDeviceCrons(devId string) {
	for _, cSt := range cr.CronByDeviceId[devId].ModulesCron {
		// ToDo: Check that 1 cSt.Cron is enough
		cSt.CronSetter(devId, cSt.Cron)
		cSt.Cron.Start()
	}
}

func (cr *DeviceCronManagerSt) ResetModuleCron(devId string, moduleName string) {
	cr.StopModuleCron(devId, moduleName)
	cr.StartModuleCron(devId, moduleName)
}

func (cr *DeviceCronManagerSt) StopModuleCron(devId string, moduleName string) {
	cr.CronByDeviceId[devId].ModulesCron[moduleName].Cron.Stop()
}

func (cr *DeviceCronManagerSt) StartModuleCron(devId string, moduleName string) {
	module := cr.CronByDeviceId[devId].ModulesCron[moduleName]
	module.CronSetter(devId, module.Cron)
	module.Cron.Start()
}

func (cr *DeviceCronManagerSt) Init() error {
	// . Go through DB and Get all Plantainers
	plantainers := PlantainersList{}
	if err := plantainerCollectionManager.FindAllModels(&bson.M{"type": "plantainer"}, &plantainers); err != nil {
		return err
	}

	fmt.Printf(time.Now().String())

	lightModuleCronSetter := func(dId string, c *cron.Cron) error {
		defer func(){
			if recover() != nil {
				return
			}
		}()
		plantainer := &PlantainerModelSt{}
		if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
			return err
		}

		lightModule := plantainer.Shadow.State.Reported.LightModule

		if *lightModule.Mode != mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
			return nil
		}

		for _, interval := range *lightModule.LightIntervalsArr {
			fromCrString := "0 " + strconv.Itoa(interval.FromTimeMinutes) + " " + strconv.Itoa(interval.FromTimeHours) + " * * *"
			c.AddFunc(fromCrString, func() {
				plantainer := &PlantainerModelSt{}
				if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
					return
				}
				lightModule := plantainer.Shadow.State.Reported.LightModule
				if *lightModule.Mode != mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
					return
				}
				plantainer.Shadow.State.Desired.LightModule.LightTurnedOn = &interval.TurnedOn
				plantainer.Shadow.IncrementVersion()

				if err := plantainerCollectionManager.SaveModel(plantainer); err != nil {
					return
				}

				successUpdate := &mccommunication.RPCMsg{
					Src: PlantainerServerId,
					Dst: dId,
					Method: dId + ".Shadow.Update.Success",
					Args: plantainer.Shadow.State,
				}
				innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", successUpdate)
				innerRPCMan.PublishRPC("User.RPC.Send", successUpdate)

				plantainer.Shadow.State.FillDelta()

				if plantainer.Shadow.State.Delta != nil {
					innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", &mccommunication.RPCMsg{
						Src: PlantainerServerId,
						Dst: dId,
						Method: dId + ".Shadow.Delta",
						Args: &map[string]interface{}{
							"state": plantainer.Shadow.State.Delta,
							"version": plantainer.Shadow.Metadata.Version,
						},
					})
				}
			})
			toCrStr := "0 " + strconv.Itoa(interval.ToTimeMinutes) + " " + strconv.Itoa(interval.ToTimeHours) + " * * *"
			c.AddFunc(toCrStr, func() {
				plantainer := &PlantainerModelSt{}
				if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
					return
				}
				lightModule := plantainer.Shadow.State.Reported.LightModule
				if *lightModule.Mode != mclightmodule.LightModuleModes[mclightmodule.LightModuleModeLightServerIntervalsTimerMode] {
					return
				}

				plantainer.Shadow.State.Desired.LightModule.LightTurnedOn = lightModule.LightIntervalsRestTimeTurnedOn
				plantainer.Shadow.IncrementVersion()

				if err := plantainerCollectionManager.SaveModel(plantainer); err != nil {
					return
				}

				successUpdate := &mccommunication.RPCMsg{
					Src: PlantainerServerId,
					Dst: dId,
					Method: dId + ".Shadow.Update.Success",
					Args: plantainer.Shadow.State,
				}
				innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", successUpdate)
				innerRPCMan.PublishRPC("User.RPC.Send", successUpdate)

				plantainer.Shadow.State.FillDelta()

				if plantainer.Shadow.State.Delta != nil {
					innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", &mccommunication.RPCMsg{
						Src: PlantainerServerId,
						Dst: dId,
						Method: dId + ".Shadow.Delta",
						Args: &map[string]interface{}{
							"state": plantainer.Shadow.State.Delta,
							"version": plantainer.Shadow.Metadata.Version,
						},
					})
				}
			})
		}
		return nil
	}

	// . Check all States and set timers
	for _, plantainer := range plantainers {
		devId := plantainer.Shadow.Id
		modulesCronMap := ModulesCronMap{}

		modulesCronMap["lightModule"] = &ModulesCronSt{
			"lightModule",
			cron.New(),
			lightModuleCronSetter,
		}

		// ToDo: Add other crons

		cr.AddDeviceCron(devId, &DeviceCronSt{
			devId,
			modulesCronMap,
			})

		cr.StartAllDeviceCrons(devId)
	}

	return nil
}

