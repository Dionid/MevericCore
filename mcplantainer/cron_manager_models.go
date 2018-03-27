package mcplantainer

import "github.com/robfig/cron"

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
