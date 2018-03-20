package mcuserrpcmanager

import (
	"mevericcore/mcws"
	"github.com/labstack/echo"
	"mevericcore/mccommon"
	"github.com/nats-io/go-nats"
)

var (
	WSManager = mcws.NewWSocketsManager()

	NATSCon *nats.Conn = nil

	UserRPCManager = CreateNewUserRPCManagerSt("plantainerServerId")

	DevicesCollectionManager mccommon.DevicesCollectionManagerInterface = nil
	UsersCollectionManager *mccommon.UsersCollectionManagerSt = nil
)

func InitMain(deviceCr DeviceCreatorFn, userColMan *mccommon.UsersCollectionManagerSt, devicesColMan mccommon.DevicesCollectionManagerInterface, e *echo.Group) {
	nc, _ := nats.Connect(nats.DefaultURL)
	NATSCon = nc

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