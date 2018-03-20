package mcdevicemqttmanager

import (
	"mevericcore/mcmqttrouter"
	"mevericcore/mccommon"
	"strings"
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
)

type DeviceMQTTManagerSt struct {
	router     *mcmqttrouter.MQTTRouter
	reqHandler mccommon.ClientToServerReqHandler
}

func (this *DeviceMQTTManagerSt) SetReqHandler(handler mccommon.ClientToServerReqHandler) {
	this.reqHandler = handler
}

func (this *DeviceMQTTManagerSt) HandleReq(c mccommon.ClientToServerHandleResChannel, msg *mccommon.ClientToServerReqSt) error {
	// ToDo: Check if this close must be somewhere else
	defer func() {
		close(c)
		if recover() != nil {
			fmt.Println("Recovered")
			return
		}
	}()

	if this.reqHandler != nil {
		return this.reqHandler(c, msg)
	}

	// ToDo: Send req to QueueManager and return

	return nil
}

func (this *DeviceMQTTManagerSt) GetChannelIdFromTopicName(topicName string) string {
	splitedTopicName := strings.Split(topicName, "/")
	return splitedTopicName[2]
}

func (this *DeviceMQTTManagerSt) GetDeviceIdFromTopicName(topicName string) string {
	splitedTopicName := strings.Split(topicName, "/")
	return splitedTopicName[4]
}

func (this *DeviceMQTTManagerSt) DeviceToServerRPCSub() {
	this.router.Subscribe("/rpc", func(client mqtt.Client, msg mqtt.Message) {
		msgPayload := msg.Payload()
		msgTopic := msg.Topic()

		fmt.Printf("Product RPC topic: %s\n", msgTopic)
		fmt.Printf("Product RPC payload: %s\n", msgPayload)

		rpcData := mccommon.RPCMsg{}
		if err := rpcData.UnmarshalJSON(msgPayload); err != nil {
			// TODO: Try to send back an error
			return
		}

		deviceId := rpcData.Src

		handleMsg := &mccommon.ClientToServerReqSt{
			ClientId:  deviceId,
			ChannelId: "",
			Protocol:  "MQTT",
			Msg:       &msgPayload,
			Resource: &msgTopic,
		}

		respChan := make(mccommon.ClientToServerHandleResChannel)

		go func() {
			err := this.HandleReq(respChan, handleMsg)
			if err != nil {
				print("OMG ERR IN MQTT CONTROLLER")
			}
		}()

		for resultSt := range respChan {
			if resultSt.Error != nil {
				if bData, err := resultSt.Error.MarshalJSON(); err != nil {
					// ToDo: Change err to RPCMsg
					this.Publish(deviceId+"/rpc", []byte(err.Error()))
				} else {
					this.Publish(deviceId+"/rpc", bData)
				}
			}
			if resultSt.Res != nil {
				if bData, err := resultSt.Res.MarshalJSON(); err != nil {
					// ToDo: Change err to RPCMsg
					this.Publish(deviceId+"/rpc", []byte(err.Error()))
				} else {
					this.Publish(deviceId+"/rpc", bData)
				}
			}
		}
	})
}

func (this *DeviceMQTTManagerSt) Publish(topic string, msg []byte) error {
	defer func() {
		if recover() != nil {
			fmt.Println("Recovered")
			return
		}
	}()

	token := this.router.Publish(topic, msg)

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return token.Error()
	}
	return nil
}

func (this *DeviceMQTTManagerSt) PublishJSON(topic string, jsonData mccommon.JSONData) error {
	defer func() {
		if recover() != nil {
			fmt.Println("Recovered")
			return
		}
	}()

	byteData, err := jsonData.MarshalJSON()
	if err != nil {
		return err
	}

	token := this.router.Publish(topic, byteData)

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return token.Error()
	}
	return nil
}

func (this *DeviceMQTTManagerSt) SendJSON(topic string, jsonData mccommon.JSONData) error {
	return this.PublishJSON(topic, jsonData)
}

func (this *DeviceMQTTManagerSt) UnSubscribeFromAll() {
	this.router.UnSubscribeFromAll()
}

func (this *DeviceMQTTManagerSt) SubscribeMainRoutes() {
	this.DeviceToServerRPCSub()
}
