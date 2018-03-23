package mcplantainer

import (
	"mevericcore/mcdevicemqttmanager"
	"mevericcore/mcmqttrouter"
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"mevericcore/mccommunication"
)

var (
	deviceMQTTMan                                        = mcdevicemqttmanager.New()
	PlantainerServerId = "plantainerServerId"
)

func initMQTT() {
	opts := mcmqttrouter.CreateConnOpts("tcp://localhost:1883", "randomString123qweasd", true)

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("OMG!!!! CONNECTION LOST BEACUSE: " + err.Error())
		deviceMQTTMan.ReInitMQTT()
	}

	mqttClient := mcmqttrouter.CreateClient(opts)

	mqttRouter := mcmqttrouter.NewMQTTRouter(mqttClient, 1)
	mqttMainG := mqttRouter.Group(PlantainerServerId)

	deviceMQTTMan.Init(mqttMainG)
	//deviceMQTTMan.SetReqHandler(deviceRPCMan.Handle)
	deviceMQTTMan.Subscribe("/rpc", func(client mqtt.Client, msg mqtt.Message) {
		msgPayload := msg.Payload()
		msgTopic := msg.Topic()

		fmt.Printf("Product RPC topic: %s\n", msgTopic)
		fmt.Printf("Product RPC payload: %s\n", msgPayload)

		rpcData := mccommunication.RPCMsg{}
		if err := rpcData.UnmarshalJSON(msgPayload); err != nil {
			// TODO: Try to send back an error
			return
		}

		deviceId := rpcData.Src

		handleMsg := &mccommunication.ClientToServerRPCReqSt{
			ClientToServerReqSt: mccommunication.ClientToServerReqSt{
				ClientId:  deviceId,
				Protocol:  "MQTT",
				Msg:       &msgPayload,
				Resource: &msgTopic,
			},
			RPCMsg: &rpcData,
		}

		respChan := make(mccommunication.ClientToServerHandleResultChannel)

		go func() {
			err := deviceRPCMan.Handle(respChan, handleMsg)
			if err != nil {
				print("OMG ERR IN MQTT CONTROLLER")
			}
		}()

		for resultSt := range respChan {
			if resultSt.Error != nil {
				if bData, err := resultSt.Error.MarshalJSON(); err != nil {
					// ToDo: Change err to RPCMsg
					deviceMQTTMan.Publish(deviceId+"/rpc", []byte(err.Error()))
				} else {
					deviceMQTTMan.Publish(deviceId+"/rpc", bData)
				}
			}
			if resultSt.Res != nil {
				if bData, err := resultSt.Res.MarshalJSON(); err != nil {
					// ToDo: Change err to RPCMsg
					deviceMQTTMan.Publish(deviceId+"/rpc", []byte(err.Error()))
				} else {
					deviceMQTTMan.Publish(deviceId+"/rpc", bData)
				}
			}
		}
	})
}
