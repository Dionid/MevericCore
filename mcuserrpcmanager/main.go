package mcuserrpcmanager

import (
	"mevericcore/mcws"
	"github.com/labstack/echo"
	"mevericcore/mccommon"
)

var (
	WSManager = mcws.NewWSocketsManager()

	UserRPCManager = CreateNewUserRPCManagerSt("plantainerServerId")

	DevicesCollectionManager mccommon.DevicesCollectionManagerInterface = nil
	UsersCollectionManager *mccommon.UsersCollectionManagerSt = nil
)

func InitMain(userColMan *mccommon.UsersCollectionManagerSt, devicesColMan mccommon.DevicesCollectionManagerInterface, e *echo.Group) {
	InitRPCManager()
	InitColManagers(userColMan, devicesColMan)
	InitHttp(e)
}

func InitRPCManager() {
	UserRPCManager.InitRoutes()
}

func InitColManagers(userColMan *mccommon.UsersCollectionManagerSt, devicesColMan mccommon.DevicesCollectionManagerInterface) {
	UsersCollectionManager = userColMan
	DevicesCollectionManager = devicesColMan
}

func InitHttp(e *echo.Group) {
	WSController := &WSHttpControllerSt{}
	e.GET("/ws", WSController.WSHandler)
}