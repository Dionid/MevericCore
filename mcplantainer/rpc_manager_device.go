package mcplantainer

import (
	"mevericcore/mclibs/mcdevicerpcmanager"
	"mevericcore/mclibs/mccommunication"
	"errors"
)

var (
	deviceRPCMan = mcdevicerpcmanager.New()
)

func initDeviceRPCManager() {
	initDeviceRPCManMainRoutes()
}

func initDeviceRPCManMainRoutes() {
	plantainerG := deviceRPCMan.Router.Group("Plantainer")
	deviceG := plantainerG.Group("Device")
	shadowG := deviceG.Group("Shadow")
	shadowG.AddHandler("Get", func(req *mccommunication.RPCReqSt) error {
		device := NewPlantainerModel()

		if err := plantainerCollectionManager.FindByShadowId(req.Msg.ClientId, device); err != nil {
			errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, "Device not found", 503)
			return deviceRPCMan.SendRPC(req.Channel, errRPC)
		}

		// ToDo: Check systems (like that intervals are working correctly)
		if changed, err := device.CheckAllSystems(); err != nil {
			errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, err.Error(), 500)
			return deviceRPCMan.SendRPC(req.Channel, errRPC)
		} else if changed {
			if err := plantainerCollectionManager.SaveModel(device); err != nil {
				errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, err.Error(), 500)
				return deviceRPCMan.SendRPC(req.Channel, errRPC)
			}
		}

		state := device.Shadow.State

		res := &map[string]interface{}{
			"state": state,
		}

		deviceRPCMan.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, res)

		state.FillDelta()

		if state.Delta != nil {
			deltaRpc := NewShadowUpdateDeltaReqRPC(device.Shadow.Id, &device.Shadow)
			deviceRPCMan.SendRPC(req.Channel, deltaRpc)
		}

		return nil
	})
	shadowG.AddHandler("Update", func(req *mccommunication.RPCReqSt) error {
		// . Get Device model and data
		device := NewPlantainerModel()

		if err := plantainerCollectionManager.FindByShadowId(req.Msg.ClientId, device); err != nil {
			errRPC := NewShadowUpdateRejectedReqRPC(req.Msg.ClientId, "Device not found", 503)
			return deviceRPCMan.SendRPC(req.Channel, errRPC)
		}

		updateRpcMsg1 := &JSONShadowUpdateRPCMsgFromDeviceSt{}

		if err := updateRpcMsg1.UnmarshalJSON(*req.Msg.Msg); err != nil {
			errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, err.Error(), 422)
			return deviceRPCMan.SendRPC(req.Channel, errRPC)
		}

		updateRpcMsg := updateRpcMsg1.ConvertToShadowUpdateRPCMsgSt()

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
					errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, err.Error(), 500)
					return deviceRPCMan.SendRPC(req.Channel, errRPC)
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
				errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, err.Error(), 500)
				return deviceRPCMan.SendRPC(req.Channel, errRPC)
			}
			device.DesiredUpdate(updateData.State.Desired)
			shadow.IncrementVersion()
		} else {
			errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, "Request is empty", 500)
			return deviceRPCMan.SendRPC(req.Channel, errRPC)
		}

		if err := plantainerCollectionManager.SaveModel(device); err != nil {
			errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, err.Error(), 500)
			return deviceRPCMan.SendRPC(req.Channel, errRPC)
		}

		// . If there was some changes to Reported, that's mean that we need to check system for ex. to change cron
		if updateData.State.Reported != nil {
			device.CheckAfterShadowReportedUpdate(&oldShadow)
		}

		// . Checking all systems to be valid and work properly
		if changed, err := device.CheckAllSystems(); err != nil {
			errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, err.Error(), 500)
			return deviceRPCMan.SendRPC(req.Channel, errRPC)
		} else if changed {
			if err := plantainerCollectionManager.SaveModel(device); err != nil {
				errRPC := NewShadowUpdateRejectedReqRPC(device.Shadow.Id, err.Error(), 500)
				return deviceRPCMan.SendRPC(req.Channel, errRPC)
			}
		}

		// . Send success back to Device
		// COMMENTED FOR A TIME
		//deviceRPCMan.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, &map[string]interface{}{
		//	"state": state,
		//	"version": device.Shadow.Metadata.Version,
		//})

		// . Send "Update.Accepted" event to Users that subscribed this device
		successUpdate := NewShadowUpdateAcceptedReqRPC(
			device.Shadow.Id,
			&device.Shadow,
		)
		deviceRPCMan.SendRPC(req.Channel, successUpdate)
		innerRPCMan.PublishRPC("User.RPC.Send", successUpdate)

		// . Check if there is some diff (delta) between Desired and Reported states (Delta struct is used for that)
		state.FillDelta()

		// . If there are some diff (delta), than send it to Device
		if state.Delta != nil {
			deltaRpc := NewShadowUpdateDeltaReqRPC(device.Shadow.Id, &device.Shadow)
			deviceRPCMan.SendRPC(req.Channel, deltaRpc)
		}

		return nil
	})
}
