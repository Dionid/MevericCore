package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
	"errors"
)

type DeviceCreatorFn func() mccommon.DeviceBaseModelInterface

type DeviceRPCCtrlSt struct{
	DeviceResponseServiceSt

	DevicesCollectionManager mccommon.DevicesCollectionManagerInterface

	Type string
	Router *DeviceRPCRouterSt

	DeviceCreator func() mccommon.DeviceBaseModelInterface
}

type DeviceRPCCtrlInterface interface {
	HandleReq(resource string, msg *mccommon.DeviceToServerReqSt, rpcData *mccommon.RPCMsg) (mccommon.JSONData, bool, mccommon.JSONData)
	SetManagers(mqttMan ProtocolManagerInterface)
}

func CreateNewDeviceRPCCtrl(typeName string, deviceCreator DeviceCreatorFn) *DeviceRPCCtrlSt {
	router := CreateNewDeviceRPCRouter()

	res := &DeviceRPCCtrlSt{
		Type: typeName,
		Router: router,
		DeviceCreator: deviceCreator,
	}

	res.InitShadowRoutes()

	return res
}

func (thisR *DeviceRPCCtrlSt) InitShadowRoutes() {
	shadowG := thisR.Router.Group("Shadow")
	shadowG.AddHandler("Get", func(req *ReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
		device := thisR.DeviceCreator()

		if err := thisR.DevicesCollectionManager.FindByShadowId(req.DeviceId, device); err != nil {
			return thisR.SendRPCErrorRes(req.Msg.Protocol, req.Msg.DeviceId, req.RPCData.Id, err.Error(), 404)
		}

		state := device.GetShadow().GetState()
		state.FillDelta()

		if len(state.Delta.State) != 0 {
			thisR.PublishDelta(req.Msg.Protocol, req.Msg.DeviceId, req.DeviceId, req.RPCData.Id, state.Delta)
		}

		return mccommon.RPCMsg{
			Src: thisR.ServerId,
			Dst: req.Msg.DeviceId,
			Id: req.RPCData.Id,
			Method: req.DeviceId + ".Shadow.Get",
			Result: &map[string]interface{}{
				"state": state,
			},
		}, true, nil
	})
	shadowG.AddHandler("Update", func(req *ReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
		device := thisR.DeviceCreator()

		if err := thisR.DevicesCollectionManager.FindByShadowId(req.DeviceId, device); err != nil {
			return thisR.SendRPCErrorRes(req.Msg.Protocol, req.Msg.DeviceId, req.RPCData.Id, err.Error(), 404)
		}

		updateRpcMsg := &mccommon.RPCWithShadowUpdateMsg{}

		if err := updateRpcMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return thisR.SendRPCErrorRes(req.Msg.Protocol, req.Msg.DeviceId, req.RPCData.Id, err.Error(), 500)
		}

		updateData := updateRpcMsg.Args

		somethingNew := false
		deviceState := device.GetShadow().GetState()

		device.ActionsOnUpdate(&updateData, thisR.DevicesCollectionManager)

		if updateData.State.Reported != nil && updateData.State.Desired != nil {
			deviceState.SetReportedState(updateData.State.Reported)
			deviceState.SetDesiredState(updateData.State.Desired)
			deviceState.IncrementVersion()
			if err := thisR.DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Msg.Protocol, req.Msg.DeviceId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Desire and Reported
			somethingNew = true
		} else if updateData.State.Reported != nil {
			deviceState.SetReportedState(updateData.State.Reported)
			deviceState.IncrementVersion()
			if err := thisR.DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Msg.Protocol, req.Msg.DeviceId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Reported
			somethingNew = true
		} else if updateData.State.Desired != nil {
			if !deviceState.CheckVersion(updateData.Version) {
				// PUB /update/rejected with Desired and Reported
				err := errors.New("version wrong")
				return thisR.SendRPCErrorRes(req.Msg.Protocol, req.Msg.DeviceId, req.RPCData.Id, err.Error(), 500)
			}
			deviceState.SetDesiredState(updateData.State.Desired)
			deviceState.IncrementVersion()
			if err := thisR.DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Msg.Protocol, req.Msg.DeviceId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Desired
			somethingNew = true
		}

		deviceState.FillDelta()

		if len(deviceState.Delta.State) != 0 {
			thisR.PublishDelta(req.Msg.Protocol, req.Msg.DeviceId, req.DeviceId, req.RPCData.Id, deviceState.Delta)
		}

		if !somethingNew {
			// In this case SetIsActivated haven't been saved
			if err := thisR.DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Msg.Protocol, req.Msg.DeviceId, req.RPCData.Id, err.Error(), 500)
			}
		}

		return mccommon.RPCMsg{
			Src: thisR.ServerId,
			Dst: req.Msg.DeviceId,
			Id: req.RPCData.Id,
			Method: req.DeviceId + ".Shadow.Update",
			Result: &map[string]interface{}{
				"state": deviceState,
			},
		}, true, nil
	})
}

func (this *DeviceRPCCtrlSt) HandleReq(resource string, msg *mccommon.DeviceToServerReqSt, rpcData *mccommon.RPCMsg) (mccommon.JSONData, bool, mccommon.JSONData) {
	return this.Router.Handle(resource, msg, rpcData)
}

func (this *DeviceRPCCtrlSt) SetManagers(mqttMan ProtocolManagerInterface) {
	this.DeviceMQTTManager = mqttMan
}


