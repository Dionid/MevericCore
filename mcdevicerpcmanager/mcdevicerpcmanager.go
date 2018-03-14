package mcdevicerpcmanager

import (
	"mevericcore/mccommon"
)

type ProtocolManagerInterface interface {
	SendJSON(string, mccommon.JSONData) error
}

type DeviceRPCManagerSt struct {
	ServerId string


	DevicesCollectionManager mccommon.DevicesCollectionManagerInterface
	//DeviceModelsAndCollectionsManager *DeviceModelsAndCollectionsManagerSt


	Router *DeviceRPCRouterSt

	DeviceHTTPManager ProtocolManagerInterface
	DeviceMQTTManager ProtocolManagerInterface
	DeviceWSManager ProtocolManagerInterface
}

func CreateDeviceRPCManager(serverId string, devicesColManager mccommon.DevicesCollectionManagerInterface, mqttMan ProtocolManagerInterface) *DeviceRPCManagerSt {
	return &DeviceRPCManagerSt{
		ServerId: serverId,
		DevicesCollectionManager: devicesColManager,
		//DeviceModelsAndCollectionsManager: CreateNewDeviceModelsAndCollectionsManager(),
		Router: CreateNewDeviceRPCRouter(),
		DeviceMQTTManager: mqttMan,
	}
}

//func (thisR *DeviceRPCManagerSt) initShadowRoutes() {
//	shadowG := thisR.Router.Group("Shadow")
//	shadowG.AddHandler("Get", func(req *ReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
//		state := req.Device.GetShadow().GetState()
//		state.FillDelta()
//
//		if len(state.Delta.State) != 0 {
//			thisR.PublishDelta(req.Msg.Protocol, req.Msg.DeviceId, req.DeviceId, req.RPCData.Id, state.Delta)
//		}
//
//		return mccommon.RPCMsg{
//			Src: thisR.ServerId,
//			Dst: req.Msg.DeviceId,
//			Id: req.RPCData.Id,
//			Result: &map[string]interface{}{
//				"state": state,
//			},
//		}, true, nil
//	})
//}

func (this *DeviceRPCManagerSt) Init() {
	//this.initShadowRoutes()
}

