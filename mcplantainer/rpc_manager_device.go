package mcplantainer

import (
	"mevericcore/mcdevicerpcmanager2"
	"mevericcore/mccommunication"
)

var (
	deviceRPCMan = mcdevicerpcmanager2.New()
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
}
