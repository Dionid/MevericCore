package mcplantainer

import (
	"mevericcore/mcdevicerpcmanager"
	"mevericcore/mccommunication"
	"mevericcore/mccommon"
)

var (
	cronRPCMan = mcdevicerpcmanager.New()
)

func initCronRPCManager() {
	initCronRPCManMainRoutes()
}

func initCronRPCManMainRoutes() {
	plantainerG := cronRPCMan.Router.Group("#")
	shadowG := plantainerG.Group("Shadow")
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
		oldShadow := device.Shadow

		// Save Data and send it to User
		if updateData.State.Reported != nil {
			data, _ := device.ExtractAndSaveData(updateData.State.Reported)
			if data != nil {
				rpcData := &mccommon.RPCMsg{
					Dst: req.Msg.RPCMsg.Src,
					Src: req.Msg.RPCMsg.Dst,
					Method: "Device.Plantainer.Data.New",
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

		// Change State
		if updateData.State.Reported != nil && updateData.State.Desired != nil {
			device.DesiredUpdate(updateData.State.Desired)
			device.ReportedUpdate(updateData.State.Reported)
			shadow.IncrementVersion()
		} else if updateData.State.Reported != nil {
			device.ReportedUpdate(updateData.State.Reported)
			shadow.IncrementVersion()
		} else if updateData.State.Desired != nil {
			device.DesiredUpdate(updateData.State.Desired)
			shadow.IncrementVersion()
		}

		if err := plantainerCollectionManager.SaveModel(device); err != nil {
			return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.Msg.RPCMsg, err.Error(), 500)
		}

		if updateData.State.Reported != nil {
			device.CheckAfterShadowReportedUpdate(&oldShadow)
		}

		deviceRPCMan.RespondSuccessResp(req.Channel, req.Msg.RPCMsg, &map[string]interface{}{
			"state": state,
			"version": device.Shadow.Metadata.Version,
		})

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

		return nil
	})
}