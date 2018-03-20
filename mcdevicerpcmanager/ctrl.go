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

	DeviceCreator DeviceCreatorFn
}

type DeviceRPCCtrlInterface interface {
	HandleReq(resource string, c mccommon.ClientToServerHandleResChannel, msg *mccommon.ClientToServerReqSt, rpcData *mccommon.RPCMsg) error
	SetManagers(mqttMan ProtocolManagerInterface)
}

func CreateNewDeviceRPCCtrl(typeName string, deviceCreator DeviceCreatorFn, mqttMan ProtocolManagerInterface) *DeviceRPCCtrlSt {
	router := CreateNewDeviceRPCRouter()

	res := &DeviceRPCCtrlSt{
		DeviceResponseServiceSt: DeviceResponseServiceSt{
			DeviceMQTTManager: mqttMan,
		},
		Type: typeName,
		Router: router,
		DeviceCreator: deviceCreator,
	}

	res.InitShadowRoutes()

	return res
}

func (thisR *DeviceRPCCtrlSt) InitShadowRoutes() {
	shadowG := thisR.Router.Group("Shadow")
	shadowG.AddHandler("Get", func(req *ReqSt) error {
		device := thisR.DeviceCreator()

		if err := thisR.DevicesCollectionManager.FindByShadowId(req.DeviceId, device); err != nil {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}

		state := device.GetShadow().GetState()
		state.FillDelta()

		if len(state.Delta.State) != 0 {
			thisR.PublishDelta(req.Channel, req.Msg.Protocol, req.Msg.ClientId, req.DeviceId, req.RPCData.Id, state.Delta)
		}

		thisR.SendRPCSuccessRes(req.Channel, req.Msg.Protocol, req.DeviceId + ".Shadow.Get", req.Msg.ClientId, req.RPCData.Id, &map[string]interface{}{
			"state": state,
		})

		return nil
	})
	shadowG.AddHandler("Update", func(req *ReqSt) error {
		device := thisR.DeviceCreator()

		if err := thisR.DevicesCollectionManager.FindByShadowId(req.DeviceId, device); err != nil {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 404)
		}

		updateRpcMsg := &mccommon.RPCWithShadowUpdateMsg{}

		if err := updateRpcMsg.UnmarshalJSON(*req.Msg.Msg); err != nil {
			return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
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
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Desire and Reported
			somethingNew = true
		} else if updateData.State.Reported != nil {
			deviceState.SetReportedState(updateData.State.Reported)
			deviceState.IncrementVersion()
			if err := thisR.DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Reported
			somethingNew = true
		} else if updateData.State.Desired != nil {
			if !deviceState.CheckVersion(updateData.Version) {
				// PUB /update/rejected with Desired and Reported
				err := errors.New("version wrong")
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
			deviceState.SetDesiredState(updateData.State.Desired)
			deviceState.IncrementVersion()
			if err := thisR.DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
			// PUB /update/accepted with Desired
			somethingNew = true
		}

		deviceState.FillDelta()

		if len(deviceState.Delta.State) != 0 {
			thisR.PublishDelta(req.Channel, req.Msg.Protocol, req.Msg.ClientId, req.DeviceId, req.RPCData.Id, deviceState.Delta)
		}

		if !somethingNew {
			// In this case SetIsActivated haven't been saved
			if err := thisR.DevicesCollectionManager.SaveModel(device); err != nil {
				return thisR.SendRPCErrorRes(req.Channel, req.Msg.Protocol, req.RPCData.Method, req.Msg.ClientId, req.RPCData.Id, err.Error(), 500)
			}
		}

		thisR.SendRPCSuccessRes(req.Channel, req.Msg.Protocol, req.DeviceId + ".Shadow.Update", req.Msg.ClientId, req.RPCData.Id, &map[string]interface{}{
			"state": deviceState,
		})

		return nil
	})
}

func (this *DeviceRPCCtrlSt) HandleReq(resource string, c mccommon.ClientToServerHandleResChannel, msg *mccommon.ClientToServerReqSt, rpcData *mccommon.RPCMsg) error {
	return this.Router.Handle(resource, c, msg, rpcData)
}

func (this *DeviceRPCCtrlSt) SetManagers(mqttMan ProtocolManagerInterface) {
	this.DeviceMQTTManager = mqttMan
}


