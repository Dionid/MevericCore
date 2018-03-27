package mcplantainer

import (
	"mevericcore/mclibs/mcuserrpcmanager"
	"mevericcore/mclibs/mccommunication"
	"mevericcore/mclibs/mccommon"
	"gopkg.in/mgo.v2/bson"
)

var (
	userRPCManager = mcuserrpcmanager.New()
)

func initUserRPCManager() {
	initUserRPCManDeviceRoutes()
}

func initUserRPCManDeviceRoutes() {
	deviceG := userRPCManager.Router.Group("Devices")
	plantainerG := deviceG.Group("Plantainer")
	plantainerG.AddHandler("Create", func(req *mccommunication.RPCReqSt) error {

		device := &PlantainerModelSt{
			DeviceBaseModel: mccommon.DeviceBaseModel{
				OwnersIds: []bson.ObjectId{bson.ObjectIdHex(req.Msg.ClientId)},
			},
			Shadow: PlantainerShadowSt{},
		}

		if err := plantainerCollectionManager.SaveModel(device); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 503)
		}

		if err := plantainerCollectionManager.FindModelById(device.ID, device); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 503)
		}

		res := &map[string]interface{}{device.Shadow.Id: device}

		return userRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
	plantainerG.AddHandler("Get", func(req *mccommunication.RPCReqSt) error {
		device := &PlantainerModelSt{}
		args := req.Msg.RPCMsg.Args.(map[string]interface{})
		deviceShadowId := args["deviceId"].(string)

		if err := plantainerCollectionManager.FindByShadowId(deviceShadowId, device); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Device not found", 503)
		}

		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 503)
		} else if !isOwner {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "You can use only your own devices", 503)
		}

		res := &map[string]interface{}{deviceShadowId: device}

		return userRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
	plantainerG.AddHandler("List", func(req *mccommunication.RPCReqSt) error {
		devices := &PlantainersList{}

		if err := plantainerCollectionManager.FindByOwnerId(req.Msg.ClientId, devices); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 503)
		}

		res := &map[string]interface{}{"data": devices}

		return userRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
	plantainerG.AddHandler("Update", func(req *mccommunication.RPCReqSt) error {
		device := &PlantainerModelSt{}
		args := req.Msg.RPCMsg.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)
		if err := plantainerCollectionManager.FindByShadowId(deviceId, device); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Device not found", 503)
		}
		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); !isOwner {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 503)
		} else if err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "You can use only your own devices", 503)
		}

		device.Update(&args)

		plantainerCollectionManager.SaveModel(device)

		res := &map[string]interface{}{deviceId: device}

		return userRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
	//shadowG := plantainerG.Group("Shadow")
}
