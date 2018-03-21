package dashboard

import (
	"github.com/labstack/echo"
	"mevericcore/mcecho"
	"mevericcore/mcuserrpcmanager"
	"mevericcore/mcplantainer/common"
	"mevericcore/mccommon"
	"mevericcore/mcusers"
	"gopkg.in/mgo.v2"
	"mevericcore/mcinnerrpc"
)

var (
	InnerRPCMan                                          = mcinnerrpc.NewInnerRPCMan()

)

func Init(session *mgo.Session, dbName string, e *echo.Echo) {
	common.InitDeviceColManager(session, dbName)

	plantainerColMan := common.PlantainerCollectionManager
	usersColMan := mccommon.InitUserColManager(session, dbName)

	initUserModules(usersColMan, e)
	initDeviceModules(usersColMan, e)

	appG := e.Group("/app")
	mcuserrpcmanager.InitMain(common.CreateNewPlantainerModelSt, common.NewPlantainersList, usersColMan, plantainerColMan, appG)

	mcuserrpcmanager.UserRPCManager.Router.AddHandler("Device.Data", func(req *mcuserrpcmanager.ReqSt) error {
		println("data in plantainer")
		device := &common.PlantainerModelSt{}
		args := req.RPCData.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)

		if err := common.PlantainerCollectionManager.FindByShadowId(deviceId, device); err != nil {
			return mcuserrpcmanager.UserRPCManager.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}
		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); !isOwner {
			return mcuserrpcmanager.UserRPCManager.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, "It's not your device", 403)
		} else if err != nil {
			print(err.Error())
			return mcuserrpcmanager.UserRPCManager.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}

		dataList := &common.PlantainerDataListSt{}

		if err := common.PlantainerDataCollectionManager.FindByDeviceShadowId(deviceId, dataList); err != nil {
			return nil
		}

		mcuserrpcmanager.UserRPCManager.SendSuccessResp(req.Channel, req.RPCData, &map[string]interface{}{"data": dataList})

		return nil
	})
}

func initUserModules(usersColMan *mccommon.UsersCollectionManagerSt, e *echo.Echo) {
	usersG := e.Group("/users")
	mcusers.InitMain(usersColMan, usersG)
}

func initDeviceModules(usersColMan *mccommon.UsersCollectionManagerSt, e *echo.Echo) {
	devicesG := e.Group("/devices")

	UserPlantainerController := &UserPlantainerControllerSt{}
	mcecho.CreateModelControllerRoutes(devicesG, "/plantainer", UserPlantainerController)
}