package mcplantainer

import (
	"mevericcore/mcdevicerpcmanager"
	"mevericcore/mccommunication"
	"mevericcore/mccommon"
	"errors"
)

var (
	deviceRPCMan = mcdevicerpcmanager.New()
)

func initDeviceRPCManager() {
	initDeviceRPCManMainRoutes()
}

func initDeviceRPCManMainRoutes() {
	plantainerG := deviceRPCMan.Router.Group("#")
	shadowG := plantainerG.Group("Shadow")
	shadowG.AddHandler("Get", func(req *mccommunication.RPCReqSt) error {
		device := NewPlantainerModel()

		if err := plantainerCollectionManager.FindByShadowId(req.Msg.ClientId, device); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Device not found", 503)
		}

		state := device.Shadow.State

		res := &map[string]interface{}{
			"state": state,
		}

		deviceRPCMan.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)

		state.FillDelta()

		if state.Delta != nil {
			deviceRPCMan.SendReq(req.Channel, req.Msg.ClientId + ".Shadow.Delta", req.Msg.RPCMsg.Dst, req.Msg.RPCMsg.Src, 123, &map[string]interface{}{
				"state": state.Delta,
				"version": device.Shadow.Metadata.Version,
			})
		}

		return nil
	})
	shadowG.AddHandler("Update", func(req *mccommunication.RPCReqSt) error {
		device := NewPlantainerModel()

		if err := plantainerCollectionManager.FindByShadowId(req.Msg.ClientId, device); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Device not found", 503)
		}

		// TODO: Can be Update result
		updateRpcMsg := &ShadowUpdateRPCMsgSt{}

		if err := updateRpcMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 422)
		}

		updateData := updateRpcMsg.Args
		shadow := &device.Shadow
		state := &device.Shadow.State

		//device.ActionsOnUpdate(&updateData, plantainerCollectionManager)

		if updateData.State.Reported != nil && updateData.State.Desired != nil {
			if state.Desired == nil {
				state.Desired = &PlantainerShadowStatePieceSt{}
			}
			state.Desired.LightModule.DesiredUpdate(&updateData.State.Desired.LightModule)
			state.Reported.LightModule.ReportedUpdate(&updateData.State.Reported.LightModule)
			shadow.IncrementVersion()
		} else if updateData.State.Reported != nil {
			//state.Reported.LightModule.ReportedUpdate(&updateData.State.Reported.LightModule)
			device.ReportedUpdate(updateData.State.Reported, plantainerCollectionManager)
			shadow.IncrementVersion()
		} else if updateData.State.Desired != nil {
			if !shadow.CheckVersion(updateData.Version) {
				err := errors.New("version wrong")
				return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
			}
			if state.Desired == nil {
				state.Desired = &PlantainerShadowStatePieceSt{}
			}
			state.Desired.LightModule.DesiredUpdate(&updateData.State.Desired.LightModule)
			shadow.IncrementVersion()
		}

		if err := plantainerCollectionManager.SaveModel(device); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
		}

		res := &map[string]interface{}{
			"state": state,
			"version": device.Shadow.Metadata.Version,
		}

		deviceRPCMan.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)

		state.FillDelta()

		if state.Delta != nil {
			deviceRPCMan.SendReq(req.Channel, req.Msg.ClientId + ".Shadow.Delta", req.Msg.RPCMsg.Dst, req.Msg.RPCMsg.Src, 123, &map[string]interface{}{
				"state":   state.Delta,
				"version": device.Shadow.Metadata.Version,
			})
		}

		rpcData := &mccommon.RPCMsg{
			Dst: req.Msg.RPCMsg.Src,
			Src: req.Msg.RPCMsg.Dst,
			Method: "Device.Shadow.Update.Accepted",
			Args: &map[string]interface{}{
				"state": updateData.State,
				"version": device.Shadow.Metadata.Version,
			},
		}

		if bData, err := rpcData.MarshalJSON(); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
		} else {
			innerRPCMan.Service.Publish("User.RPC.Send", bData)
		}

		return nil
	})
}
