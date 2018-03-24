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
		state.FillDelta()

		if len(state.Delta.State) != 0 {
			deviceRPCMan.SendReq(req.Channel, req.Msg.ClientId + ".Shadow.Delta", req.Msg.RPCMsg.Dst, req.Msg.RPCMsg.Src, 123, &map[string]interface{}{
				"state": state.Delta.State,
				"version": state.Delta.Version,
			})
		}

		res := &map[string]interface{}{
			"state": state,
		}

		return deviceRPCMan.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)
	})
	shadowG.AddHandler("Update", func(req *mccommunication.RPCReqSt) error {
		device := NewPlantainerModel()

		if err := plantainerCollectionManager.FindByShadowId(req.Msg.ClientId, device); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, "Device not found", 503)
		}

		// TODO: Can be Update result

		updateRpcMsg := &mccommon.RPCWithShadowUpdateMsg{}

		if err := updateRpcMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 422)
		}

		updateData := updateRpcMsg.Args

		somethingNew := false
		deviceState := device.GetShadow().GetState()

		device.ActionsOnUpdate(&updateData, plantainerCollectionManager)

		if updateData.State.Reported != nil && updateData.State.Desired != nil {
			deviceState.SetReportedState(updateData.State.Reported)
			deviceState.SetDesiredState(updateData.State.Desired)
			deviceState.IncrementVersion()
			if err := plantainerCollectionManager.SaveModel(device); err != nil {
				return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
			}
			// PUB /update/accepted with Desire and Reported
			somethingNew = true
		} else if updateData.State.Reported != nil {
			deviceState.SetReportedState(updateData.State.Reported)
			deviceState.IncrementVersion()
			if err := plantainerCollectionManager.SaveModel(device); err != nil {
				return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
			}
			// PUB /update/accepted with Reported
			somethingNew = true
		} else if updateData.State.Desired != nil {
			if !deviceState.CheckVersion(updateData.Version) {
				// PUB /update/rejected with Desired and Reported
				err := errors.New("version wrong")
				return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
			}
			deviceState.SetDesiredState(updateData.State.Desired)
			deviceState.IncrementVersion()
			if err := plantainerCollectionManager.SaveModel(device); err != nil {
				return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
			}
			// PUB /update/accepted with Desired
			somethingNew = true
		}

		deviceState.FillDelta()

		if len(deviceState.Delta.State) != 0 {
			deviceRPCMan.SendReq(req.Channel, req.Msg.ClientId + ".Shadow.Delta", req.Msg.RPCMsg.Dst, req.Msg.RPCMsg.Src, 123, &map[string]interface{}{
				"state": deviceState.Delta.State,
				"version": deviceState.Delta.Version,
			})
		}

		if !somethingNew {
			// In this case SetIsActivated haven't been saved
			if err := plantainerCollectionManager.SaveModel(device); err != nil {
				return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
			}
		}

		res := &map[string]interface{}{
			"state": deviceState,
		}

		deviceRPCMan.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)

		rpcData := &mccommon.RPCMsg{
			Dst: req.Msg.RPCMsg.Src,
			Src: req.Msg.RPCMsg.Dst,
			Method: "Device.Shadow.Update.Accepted",
			Args: &map[string]interface{}{
				"state": updateData.State,
				"version": deviceState.Metadata.Version,
			},
		}

		if bData, err := rpcData.MarshalJSON(); err != nil {
			//return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
		} else {
			innerRPCMan.Service.Publish("User.RPC.Send", bData)
		}

		return nil
	})
}
