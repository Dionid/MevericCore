package mcplantainer

import (
	"mevericcore/mcdevicerpcmanager"
	"mevericcore/mccommunication"
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

		// ToDo: Check systems (like that intervals are working correctly)

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
		// . Get Device model and data
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
		oldShadow := device.Shadow

		// . Save if there is some incoming Data to store (like LightLvl or Humidity)
		if updateData.State.Reported != nil {
			data, _ := device.ExtractAndSaveData(updateData.State.Reported)
			if data != nil {
				// . Send new Data to User
				rpcData := &mccommunication.RPCMsg{
					Dst: req.Msg.RPCMsg.Src,
					Src: req.Msg.RPCMsg.Dst,
					Method: "Plantainer.Device.Data.New",
					Args: &map[string]interface{}{
						"data": data,
					},
				}

				if bData, err := rpcData.MarshalJSON(); err != nil {
					return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
				} else {
					innerRPCMan.Service.Publish("User.RPC.Send", bData)
				}
			}
		}

		// . Change Device.Shadow.State
		if updateData.State.Reported != nil && updateData.State.Desired != nil {
			device.DesiredUpdate(updateData.State.Desired)
			device.ReportedUpdate(updateData.State.Reported)
			shadow.IncrementVersion()
		} else if updateData.State.Reported != nil {
			device.ReportedUpdate(updateData.State.Reported)
			shadow.IncrementVersion()
		} else if updateData.State.Desired != nil {
			if !shadow.CheckVersion(updateData.Version) {
				err := errors.New("version wrong")
				return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
			}
			device.DesiredUpdate(updateData.State.Desired)
			shadow.IncrementVersion()
		}

		if err := plantainerCollectionManager.SaveModel(device); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
		}

		// . If there was some changes to Reported, that's mean that we need to check system for ex. to change cron
		if updateData.State.Reported != nil {
			device.CheckAfterShadowReportedUpdate(&oldShadow)
		}

		// . Send success back to Device
		deviceRPCMan.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, &map[string]interface{}{
			"state": state,
			"version": device.Shadow.Metadata.Version,
		})

		// . Send "Update.Accepted" event to Users that subscribed this device
		rpcData := &mccommunication.RPCMsg{
			Dst: req.Msg.RPCMsg.Src,
			Src: req.Msg.RPCMsg.Dst,
			Method: "Device.Plantainer.Shadow.Update.Accepted",
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

		// . Check if there is some diff (delta) between Desired and Reported states (Delta struct is used for that)
		state.FillDelta()

		// . If there are some diff (delta), than send it to Device
		if state.Delta != nil {
			deviceRPCMan.SendReq(req.Channel, req.Msg.ClientId + ".Shadow.Delta", req.Msg.RPCMsg.Dst, req.Msg.RPCMsg.Src, 123, &map[string]interface{}{
				"state":   state.Delta,
				"version": device.Shadow.Metadata.Version,
			})
		}

		return nil
	})
}
