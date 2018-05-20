package mcplantainer

import (
	"mevericcore/mclibs/mcdevicemqttmanager"
	"mevericcore/mclibs/mcmqttrouter"
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"mevericcore/mclibs/mccommunication"
)

var (
	deviceMQTTMan                                        = mcdevicemqttmanager.New()
	PlantainerServerId = "plantainerServerId"
)

func subscribeToRPC() {
	fmt.Println("")
	fmt.Println("Subscribing to main routes")
	fmt.Println("")
	deviceMQTTMan.Subscribe("/rpc", func(client mqtt.Client, msg mqtt.Message) {
		msgPayload := msg.Payload()
		msgTopic := msg.Topic()

		fmt.Printf("Plantainer received RPC topic: %s\n", msgTopic)
		fmt.Printf("Plantainer received RPC payload: %s\n", msgPayload)

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
				println("OMG ERR IN MQTT CONTROLLER: " + err.Error())
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

func initMQTT() {
	//opts := mcmqttrouter.CreateConnOpts("tcp://195.201.91.242:1883", "randomString123qweasd", true)
	opts := mcmqttrouter.CreateConnOpts("tcp://iot.eclipse.org:1883", "randomString123qweasd", false, true, 30, 30)
	//opts := mcmqttrouter.CreateConnOpts("tcp://localhost:1883", "randomString123qweasd", true)

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("OMG!!!! CONNECTION LOST BEACUSE: " + err.Error())
		deviceMQTTMan.ReInitMQTT()
		subscribeToRPC()
	}

	mqttClient := mcmqttrouter.CreateClient(opts)

	mqttRouter := mcmqttrouter.NewMQTTRouter(mqttClient, 1)
	mqttMainG := mqttRouter.Group(PlantainerServerId)

	deviceMQTTMan.Init(mqttMainG)
	//deviceMQTTMan.SetReqHandler(deviceRPCMan.Handle)
	subscribeToRPC()
}
