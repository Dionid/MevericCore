package mcplantainer

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"mevericcore/mcmqttrouter"
	"mevericcore/mcdevicemqttmanager"
	"gopkg.in/mgo.v2"
	"mevericcore/mcdevicerpcmanager"
)



var (
	DeviceRPCManager = mcdevicerpcmanager.CreateDeviceRPCManager("plantainerServerId", &DevicesCollectionManager, mcdevicemqttmanager.DeviceMQTTManager)
)

func activateMQTT() {
	//opts := mcmqttrouter.CreateConnOpts("tcp://iot.eclipse.org:1883", "randomString123qweasd", true)
	opts := mcmqttrouter.CreateConnOpts("tcp://localhost:1883", "randomString123qweasd", true)
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("OMG!!!! CONNECTION LOST BEACUSE: " + err.Error())
		mcdevicemqttmanager.ReInitMQTT()
	}
	c := mcmqttrouter.CreateClient(opts)
	mqttRouter := mcmqttrouter.NewMQTTRouter(c)
	mqttMainG := mqttRouter.Group("plantainerServerId")

	mcdevicemqttmanager.Init(mqttMainG)
	mcdevicemqttmanager.DeviceMQTTManager.SetReqHandler(DeviceRPCManager.RPCReqHandler)

	fmt.Println("MQTT IS ACTIVATED")
}

func Init(dbsession *mgo.Session, dbName string) {
	initDeviceColManager(dbsession, dbName)

	DeviceRPCManager.DeviceModelsAndCollectionsManager.RegisterNewDeviceType("plantainer", CreateNewPlantainerModelSt, &DevicesCollectionManager)

	// 1. Activate MQTT
	activateMQTT()

	// 2. Activate HTTP

}