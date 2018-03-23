package mcplantainer

import (
	"mevericcore/mcuserrpcmanager2"
	"mevericcore/mccommunication"
)

var (
	UserRPCManager = mcuserrpcmanager2.New()
)

func InitUserRPCManager() {
	initUserRPCManDeviceRoutes()
}

func initUserRPCManDeviceRoutes() {
	deviceG := UserRPCManager.Router.Group("Devices")
	plantainerG := deviceG.Group("Plantainer")
	plantainerG.AddHandler("Get", func(req *mccommunication.RPCReqSt) error {
		return nil
	})
}
