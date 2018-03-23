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
		device := &PlantainerModelSt{}
		args := req.Msg.RPCMsg.Args.(map[string]interface{})
		deviceShadowId := args["deviceId"].(string)

		if err := plantainerCollectionManager.FindByShadowId(deviceShadowId, device); err != nil {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Device not found", 503)
		}

		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); err != nil {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 503)
		} else if !isOwner {
			return UserRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "You can use only your own devices", 503)
		}

		res := &map[string]interface{}{deviceShadowId: device}

		return UserRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
}
