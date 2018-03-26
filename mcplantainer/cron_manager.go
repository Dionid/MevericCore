package mcplantainer

import (
	"github.com/robfig/cron"
	"mevericcore/mcinnerrpc"
	"mevericcore/mccommunication"
)

type CronSetterFn func(devId string, lightModuleC *cron.Cron) error

type ModulesCronSt struct {
	ModuleName string
	Cron *cron.Cron
	CronSetter CronSetterFn
}

type ModulesCronMap map[string]*ModulesCronSt

type DeviceCronSt struct {
	DeviceShadowId string
	ModulesCron ModulesCronMap
}

type DeviceCronInterface interface {

}

// This manager is used in Hub to check operations that needs to be done
type DeviceCronManager struct {
	CronByDeviceId map[string]*DeviceCronSt
}

func (cr *DeviceCronManager) AddCron(devId string, c *DeviceCronSt) {
	cr.CronByDeviceId[devId] = c
}

func (cr *DeviceCronManager) StartAllDeviceCrons(devId string) {
	for _, cSt := range cr.CronByDeviceId[devId].ModulesCron {
		cSt.CronSetter(devId, cSt.Cron)
		cSt.Cron.Start()
	}
}

func (cr *DeviceCronManager) ResetModuleCron(devId string, moduleName string) {
	cr.StopModuleCron(devId, moduleName)
	cr.StartModuleCron(devId, moduleName)
}

func (cr *DeviceCronManager) StopModuleCron(devId string, moduleName string) {
	cr.CronByDeviceId[devId].ModulesCron[moduleName].Cron.Stop()
}

func (cr *DeviceCronManager) StartModuleCron(devId string, moduleName string) {
	module := cr.CronByDeviceId[devId].ModulesCron[moduleName]
	module.CronSetter(devId, module.Cron)
}

func (cr *DeviceCronManager) subInnerRPC() {
	// . Subscribe for tasks
	innerRPCMan.Service.Subscribe("DeviceCron.Plantainer.RPC", func(msg *mcinnerrpc.Msg) {
		rpcData := mccommunication.RPCMsg{}

		if err := rpcData.UnmarshalJSON(msg.Data); err != nil {
			return
		}

		args := rpcData.Args.(map[string]interface{})
		devId := args["deviceId"].(string)
		modules := args["modules"].([]string)

		switch rpcData.Method {
		case "DeviceCron.Plantainer.RPC.Stop":
			for _, name := range modules {
				cr.StopModuleCron(devId, name)
			}
		case "DeviceCron.Plantainer.RPC.Reset":
			for _, name := range modules {
				cr.ResetModuleCron(devId, name)
			}
		}
	})
}

func (cr *DeviceCronManager) Init() error {
	cr.subInnerRPC()

	// . Go through DB and Get all Plantainers
	plantainers := PlantainersList{}
	if err := plantainerCollectionManager.FindAllModels(nil, &plantainers); err != nil {
		return err
	}

	// . Check all States and set timers
	for _, plantainer := range plantainers {
		devId := plantainer.Shadow.Id
		modulesCronMap := ModulesCronMap{}

		lightModuleCronSetter := func(dId string, c *cron.Cron) error {
			plantainer := &PlantainerModelSt{}
			if err := plantainerCollectionManager.FindByShadowId(dId, plantainer); err != nil {
				return err
			}
			plantainer.Shadow.State.Reported.LightModule.SetCronTasks(dId, c)
			return nil
		}

		modulesCronMap["lightModule"] = &ModulesCronSt{
			"lightModule",
			cron.New(),
			lightModuleCronSetter,
		}

		cr.AddCron(devId, &DeviceCronSt{
			devId,
			modulesCronMap,
			})

		cr.StartAllDeviceCrons(devId)
	}

	return nil
}

