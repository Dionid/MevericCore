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

func InitMain(deviceCr DeviceCreatorFn, userColMan *mccommon.UsersCollectionManagerSt, devicesColMan mccommon.DevicesCollectionManagerInterface, e *echo.Group) {
	InitRPCManager(deviceCr)
	InitColManagers(userColMan, devicesColMan)
	InitHttp(e)
}

func InitRPCManager(deviceCr DeviceCreatorFn) {
	UserRPCManager.Init(deviceCr)
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