package mcdevicerpcmanager

import "mevericcore/mccommon"

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
			Result: &map[string]interface{}{
				"state": state,
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


