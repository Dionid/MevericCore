package mcplantainer

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"mevericcore/mcmqttrouter"
	"mevericcore/mcdevicemqttmanager"
	"mevericcore/mccommon"
	"strings"
	"gopkg.in/mgo.v2"
)

func DeviceRPCReqHandler (msg *mccommon.DeviceToServerReqSt) (res mccommon.JSONData, sendBack bool, err mccommon.JSONData) {
	rpcData := mcdevicemqttmanager.RPCMsg{}
	if err := rpcData.UnmarshalJSON(*msg.Msg); err != nil {
		return nil, true, mcdevicemqttmanager.RPCMsg{
			Src: "qweyuiasdhjky",
			Dst: msg.DeviceId,
			Id: rpcData.Id,
			Error: &map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	splitedMethod := strings.Split(rpcData.Method, ".")

	if splitedMethod[0] != msg.DeviceId {
		return nil, false, nil
	}

	if splitedMethod[1] == "Shadow" {
		if splitedMethod[2] == "Get" {
			plantainer := PlantainerModelSt{}

			// 1. Find device model
			if err := DevicesCollectionManager.FindByShadowId(msg.DeviceId, &plantainer); err != nil {
				return nil, true, mcdevicemqttmanager.RPCMsg{
					Src: "qweyuiasdhjky",
					Dst: msg.DeviceId,
					Id: rpcData.Id,
					Error: &map[string]interface{}{
						"error": err.Error(),
					},
				}
			}

			// 2. Send state in Get.Accepted
			state := plantainer.Shadow.State

			state.FillDelta()

			if state.Delta != nil {
				mcdevicemqttmanager.DeviceMQTTManager.PublishJSON(msg.DeviceId+"/rpc", mcdevicemqttmanager.RPCMsg{
					Src: "qweyuiasdhjky",
					Dst: msg.DeviceId,
					Id: rpcData.Id,
					Method: msg.DeviceId+".Shadow.Delta",
					Args: state.Delta,
				})
			}
		}
		if splitedMethod[2] == "Update" {

		}
	}

	return mcdevicemqttmanager.RPCMsg{
		Src: "qweyuiasdhjky",
		Dst: msg.DeviceId,
		Id: rpcData.Id,
		Result: &map[string]interface{}{
			"success": true,
		},
	}, true, nil
}

func activateMQTT() {
	//opts := mcmqttrouter.CreateConnOpts("tcp://iot.eclipse.org:1883", "randomString123qweasd", true)
	opts := mcmqttrouter.CreateConnOpts("tcp://localhost:1883", "randomString123qweasd", true)
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("OMG!!!! CONNECTION LOST BEACUSE: " + err.Error())
		mcdevicemqttmanager.ReInitMQTT()
	}
	c := mcmqttrouter.CreateClient(opts)
	mqttRouter := mcmqttrouter.NewMQTTRouter(c)
	mqttMainG := mqttRouter.Group("qweyuiasdhjky")

	mcdevicemqttmanager.Init(mqttMainG)
	mcdevicemqttmanager.DeviceMQTTManager.SetReqHandler(DeviceRPCReqHandler)

	fmt.Println("MQTT IS ACTIVATED")
}

func Init(dbsession *mgo.Session, dbName string) {
	initDeviceColManager(dbsession, dbName)

	// 1. Activate MQTT
	activateMQTT()

	// 2. Activate HTTP
}