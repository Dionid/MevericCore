package mcplantainer

import (
	"mevericcore/mclibs/mcuserrpcmanager"
	"mevericcore/mclibs/mccommunication"
	"mevericcore/mclibs/mccommon"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

var (
	userRPCManager = mcuserrpcmanager.New()
)

func initUserRPCManager() {
	initUserRPCManDeviceRoutes()
}

func initUserRPCManDeviceRoutes() {
	plantainerG := userRPCManager.Router.Group("Plantainer")
	deviceG := plantainerG.Group("Devices")
	deviceG.AddHandler("Create", func(req *mccommunication.RPCReqSt) error {

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
	deviceG.AddHandler("Get", func(req *mccommunication.RPCReqSt) error {
		device := &PlantainerModelSt{}
		args := req.Msg.RPCMsg.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)

		if err := plantainerCollectionManager.FindByShadowId(deviceId, device); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Device not found", 503)
		}

		if isOwner, err := device.IsOwnerStringId(req.Msg.ClientId); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 503)
		} else if !isOwner {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "You can use only your own devices", 503)
		}

		res := &map[string]interface{}{deviceId: device}

		device.Shadow.State.FillDelta()

		// . If there are some diff (delta), than send it to Device
		if device.Shadow.State.Delta != nil {
			rpcData := NewShadowUpdateDeltaReqRPC(deviceId, &device.Shadow)
			//rpcData := &mccommunication.RPCMsg{
			//	Method: deviceId + ".Shadow.Delta",
			//	Id: req.Msg.RPCMsg.Id,
			//	Src: PlantainerServerId,
			//	Dst: deviceId,
			//	Args: &map[string]interface{}{
			//		"state":   device.Shadow.State.Delta,
			//		"version": device.Shadow.Metadata.Version,
			//	},
			//}
			innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", rpcData)
		}

		return userRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
	deviceG.AddHandler("List", func(req *mccommunication.RPCReqSt) error {
		devices := &PlantainersList{}

		if err := plantainerCollectionManager.FindByOwnerId(req.Msg.ClientId, devices); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Try again", 503)
		}

		res := &map[string]interface{}{"data": devices}

		return userRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
	deviceG.AddHandler("Update", func(req *mccommunication.RPCReqSt) error {
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
	shadowG := deviceG.Group("Shadow")
	shadowG.AddHandler("Update", func(req *mccommunication.RPCReqSt) error {
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

		updateRpcMsg := &ShadowUpdateRPCMsgSt{}

		if err := updateRpcMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 422)
		}

		updateData := updateRpcMsg.Args
		shadow := &device.Shadow
		state := &device.Shadow.State
		//oldShadow := device.Shadow

		if updateData.State.Desired != nil {
			if !shadow.CheckVersion(updateData.Version) {
				err := errors.New("version wrong")
				return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
			}
			device.DesiredUpdate(updateData.State.Desired)
			shadow.IncrementVersion()
		} else {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Shadow.Desired is required", 500)
		}

		if err := plantainerCollectionManager.SaveModel(device); err != nil {
			return userRPCManager.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
		}

		// . Send success back to Device
		userRPCManager.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, &map[string]interface{}{
			"state": state,
			"version": device.Shadow.Metadata.Version,
		})

		successUpdate := NewShadowUpdateAcceptedReqRPC(
			deviceId,
			&device.Shadow,
		)
		deviceRPCMan.SendRPC(req.Channel, successUpdate)

		state.FillDelta()

		// . If there are some diff (delta), than send it to Device
		if state.Delta != nil {
			//rpcData := &mccommunication.RPCMsg{
			//	Method: deviceId + ".Shadow.Delta",
			//	Id: req.Msg.RPCMsg.Id,
			//	Src: PlantainerServerId,
			//	Dst: deviceId,
			//	Args: &map[string]interface{}{
			//		"state":   state.Delta,
			//		"version": device.Shadow.Metadata.Version,
			//	},
			//}
			deltaRpc := NewShadowUpdateDeltaReqRPC(deviceId, &device.Shadow)
			innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", deltaRpc)
		}

		return nil
	})
}
