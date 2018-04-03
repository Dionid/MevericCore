package mcplantainer

import (
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"time"
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

	// . Check all States and set timers
	for _, plantainer := range plantainers {
		devId := plantainer.Shadow.Id
		modulesCronMap := ModulesCronMap{}
		modulesCronMap["lightModule"] = &ModulesCronSt{
			"lightModule",
			cron.New(),
			cr.NewLightModuleSetter(),
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

