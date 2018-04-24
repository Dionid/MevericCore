package mcplantainer

import (
	"mevericcore/mclibs/mcrpcmanager"
	"mevericcore/mclibs/mcrpcrouter"
)

var (
	cronRPCMan = mcrpcmanager.New()
)

func initCronRPCManager() {
	initCronRPCManMainRoutes()
}

func initCronRPCManMainRoutes() {
	plantainerG := cronRPCMan.Router.Group("Plantainer")
	cronG := plantainerG.Group("Cron")
	cronG.AddHandler("Reset", func(req *mcrpcrouter.RPCReqSt) error {
		args := req.RPCData.Args.(map[string]interface{})
		devId := args["deviceId"].(string)
		modules := args["modules"].([]interface{})
		for _, name := range modules {
			deviceCronManager.ResetModuleCron(devId, name.(string))
		}
		return nil
	})
	cronG.AddHandler("Stop", func(req *mcrpcrouter.RPCReqSt) error {
		args := req.RPCData.Args.(map[string]interface{})
		devId := args["deviceId"].(string)
		modules := args["modules"].([]interface{})
		for _, name := range modules {
			deviceCronManager.StopModuleCron(devId, name.(string))
		}
		return nil
	})
}
