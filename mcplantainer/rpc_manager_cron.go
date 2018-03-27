package mcplantainer

import (
	"mevericcore/mccommunication"
	"mevericcore/mclibs/mcrpcmanager"
	"mevericcore/mclibs/mcrpcrouter"
	"errors"
)

var (
	cronRPCMan = mcrpcmanager.New()
)

func initCronRPCManager() {
	initCronRPCManMainRoutes()
}

func initCronRPCManMainRoutes() {
	plantainerG := cronRPCMan.Router.Group("Plantainer")
	cronG := plantainerG.Group("Cron")
	cronG.AddHandler("Reset", func(req *mcrpcrouter.RPCReqSt) error {
		args := req.RPCData.Args.(map[string]interface{})
		devId := args["deviceId"].(string)
		modules := args["modules"].([]interface{})
		for _, name := range modules {
			deviceCronManager.ResetModuleCron(devId, name.(string))
		}
		return nil
	})
	cronG.AddHandler("Stop", func(req *mcrpcrouter.RPCReqSt) error {
		args := req.RPCData.Args.(map[string]interface{})
		devId := args["deviceId"].(string)
		modules := args["modules"].([]interface{})
		for _, name := range modules {
			deviceCronManager.StopModuleCron(devId, name.(string))
		}
		return nil
	})

	deviceG := plantainerG.Group("Device")
	shadowG := deviceG.Group("Shadow")
	shadowG.AddHandler("Update", func(req *mcrpcrouter.RPCReqSt) error {
		device := NewPlantainerModel()

		args := req.RPCData.Args.(map[string]interface{})
		deviceId := args["deviceId"].(string)

		if err := plantainerCollectionManager.FindByShadowId(deviceId, device); err != nil {
			return cronRPCMan.RespondRPCErrorRes(req.Channel, req.RPCData, "Device not found", 503)
		}

		updateRpcMsg := &ShadowUpdateRPCMsgSt{}

		if err := updateRpcMsg.UnmarshalJSON(*req.Msg); err != nil {
			return cronRPCMan.RespondRPCErrorRes(req.Channel, req.RPCData, err.Error(), 422)
		}

		updateData := updateRpcMsg.Args
		shadow := &device.Shadow
		state := &device.Shadow.State
		oldShadow := device.Shadow

		// Change State
		if updateData.State.Desired != nil {
			if !shadow.CheckVersion(updateData.Version) {
				err := errors.New("version wrong")
				return deviceRPCMan.RespondRPCErrorRes(req.Channel, req.RPCData, err.Error(), 500)
			}
			device.DesiredUpdate(updateData.State.Desired)
			shadow.IncrementVersion()
		}

		if err := plantainerCollectionManager.SaveModel(device); err != nil {
			return cronRPCMan.RespondRPCErrorRes(req.Channel, req.RPCData, err.Error(), 500)
		}

		if updateData.State.Reported != nil {
			// ToDo: Check that this function is needed right here
			device.CheckAfterShadowReportedUpdate(&oldShadow)
		}

		//cronRPCMan.RespondSuccessResp(req.Channel, req.RPCData, &map[string]interface{}{
		//	"state": state,
		//	"version": device.Shadow.Metadata.Version,
		//})

		rpcData := &mccommunication.RPCMsg{
			Dst: deviceId,
			Src: PlantainerServerId,
			Method: "Plantainer.Shadow.Update.Accepted",
			Args: &map[string]interface{}{
				"state": updateData.State,
				"version": device.Shadow.Metadata.Version,
			},
		}

		innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", rpcData)
		innerRPCMan.PublishRPC("User.RPC.Send", rpcData)

		state.FillDelta()

		if state.Delta != nil {
			rpcData := &mccommunication.RPCMsg{
				Dst: deviceId,
				Src: PlantainerServerId,
				Method: deviceId + ".Shadow.Delta",
				Args: &map[string]interface{}{
					"state":   state.Delta,
					"version": device.Shadow.Metadata.Version,
				},
			}
			innerRPCMan.PublishRPC("Plantainer.Device.RPC.Send", rpcData)
		}

		return nil
	})
}
