package mcuserrpcmanager

import (
	"mevericcore/mcws"
	"github.com/labstack/echo"
	"mevericcore/mccommon"
	"mevericcore/mcinnerrpc"
)

var (
	WSManager = mcws.NewWSocketsManager()

	InnerRPCMan                                          = mcinnerrpc.New()

	// ToDo: "plantainerServerId" can't be here
	UserRPCManager = CreateNewUserRPCManagerSt("plantainerServerId")

	DevicesCollectionManager mccommon.DevicesCollectionManagerInterface = nil
	UsersCollectionManager *mccommon.UsersCollectionManagerSt = nil
)

func InitMain(deviceCr mccommon.DeviceCreatorFn, devicesLCr mccommon.DevicesListCreatorFn, userColMan *mccommon.UsersCollectionManagerSt, devicesColMan mccommon.DevicesCollectionManagerInterface, e *echo.Group) {
	InitInnerRPCManager()

	InitRPCManager(deviceCr, devicesLCr)
	InitColManagers(userColMan, devicesColMan)
	InitHttp(e)
}

func InitInnerRPCManager() {
	InnerRPCMan.Init()
	InnerRPCMan.Service.Subscribe("User.RPC.Send", func(msg *mcinnerrpc.Msg) {
		rpcData := &mcws.WsRPCMsgBaseSt{}

		if err := rpcData.UnmarshalJSON(msg.Data); err != nil {
			return
		}

		WSManager.SendWsMsgByRoomName(rpcData.Dst, rpcData)
	})
}

func InitColManagers(userColMan *mccommon.UsersCollectionManagerSt, devicesColMan mccommon.DevicesCollectionManagerInterface) {
	UsersCollectionManager = userColMan
	DevicesCollectionManager = devicesColMan
}

func InitRPCManager(deviceCr mccommon.DeviceCreatorFn, devicesLCr mccommon.DevicesListCreatorFn) {
	UserRPCManager.Init(deviceCr, devicesLCr)
	UserRPCManager.InitRoutes()
}

func InitHttp(e *echo.Group) {
	WSController := &WSHttpControllerSt{}
	e.GET("/ws", WSController.WSHandler)
}