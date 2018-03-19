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
	reqHandler mccommon.DeviceToServerReqHandler
	//deviceIdPropName string
}

func (this *DeviceMQTTManagerSt) SetReqHandler(handler mccommon.DeviceToServerReqHandler) {
	this.reqHandler = handler
}

func (this *DeviceMQTTManagerSt) HandleReq(msg *mccommon.DeviceToServerReqSt) (resMsg mccommon.JSONData, sendBack bool, errMsg mccommon.JSONData) {
	if this.reqHandler != nil {
		return this.reqHandler(msg)
	}

	// ToDo: Send req to QueueManager and return

	return nil, false, nil
}

func (this *DeviceMQTTManagerSt) GetChannelIdFromTopicName(topicName string) string {
	splitedTopicName := strings.Split(topicName, "/")
	return splitedTopicName[2]
}

func (this *DeviceMQTTManagerSt) GetDeviceIdFromTopicName(topicName string) string {
	splitedTopicName := strings.Split(topicName, "/")
	return splitedTopicName[4]
}

func (this *DeviceMQTTManagerSt) DeviceToServerSub() {
	this.router.Subscribe("channels/+/devices/+/s/rpc", func(client mqtt.Client, msg mqtt.Message) {
		msgPayload := msg.Payload()
		msgTopic := msg.Topic()

		fmt.Printf("Product RPC topic: %s\n", msgTopic)
		fmt.Printf("Product RPC payload: %s\n", msgPayload)

		deviceId := this.GetDeviceIdFromTopicName(msgTopic)
		channelId := this.GetChannelIdFromTopicName(msgTopic)

		// ToDo: Auth
		// Get PublisherId
		// MUST: deviceId == publisherId

		// ToDo: Replace this code to other module
		//device := new(DeviceBaseModelSt)
		//if err := DevicesCollectionManager.FindByStringId(deviceId, device); err != nil {
		//	return
		//}
		//
		//channel := new(ChannelBaseModel)
		//if err := ChannelsCollectionManager.FindByStringId(channelId, channel); err != nil {
		//	return
		//}
		//
		//if err := channel.CheckPublisherDeviceAccess(publisherId, device); err != nil {
		//	return
		//}
		//
		//man := GlobalActionManager.DeviceToServerActionManagersByDeviceType[device.Type]
		//
		//res, sendBack, err := man.ProcessAction()

		res, sendBack, err := this.HandleReq(&mccommon.DeviceToServerReqSt{
			DeviceId:  deviceId,
			ChannelId: channelId,
			Protocol:  "MQTT",
			Msg:       &msgPayload,
		})

		if err != nil {
			// LOG
			b, jerr := res.MarshalJSON()
			if jerr != nil {
				// LOG
				return
			}
			// TODO: Change this to normal way
			this.Publish(deviceId+"/rpc", b)
		}

		if sendBack && res != nil {
			b, jerr := res.MarshalJSON()
			if jerr != nil {
				// LOG
				return
			}
			this.Publish(deviceId+"/rpc", b)
		}
	})
}

func (this *DeviceMQTTManagerSt) DevicetoServerRPCSub() {
	this.router.Subscribe("/rpc", func(client mqtt.Client, msg mqtt.Message) {
		msgPayload := msg.Payload()
		msgTopic := msg.Topic()

		fmt.Printf("Product RPC topic: %s\n", msgTopic)
		fmt.Printf("Product RPC payload: %s\n", msgPayload)

		rpcData := mccommon.RPCMsg{}
		if err := rpcData.UnmarshalJSON(msgPayload); err != nil {
			//b, jerr := res.MarshalJSON()
			//if jerr != nil {
			//	// LOG
			//	return
			//}
			//this.Publish(deviceId+"/rpc", b)
			return
		}

		deviceId := rpcData.Src

		res, sendBack, err := this.HandleReq(&mccommon.DeviceToServerReqSt{
			DeviceId:  deviceId,
			ChannelId: "",
			Protocol:  "MQTT",
			Msg:       &msgPayload,
		})

		if err != nil {
			// LOG
			b, jerr := err.MarshalJSON()
			if jerr != nil {
				// LOG
				return
			}
			this.Publish(deviceId+"/rpc", b)
			return
		}

		if sendBack && res != nil {
			b, jerr := res.MarshalJSON()
			if jerr != nil {
				// LOG
				return
			}
			this.Publish(deviceId+"/rpc", b)
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
	this.DevicetoServerRPCSub()
}
