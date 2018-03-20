package device

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"mevericcore/mcmqttrouter"
	"mevericcore/mcdevicemqttmanager"
	"gopkg.in/mgo.v2"
	"mevericcore/mcdevicerpcmanager"
	"mevericcore/mcplantainer/common"
)

var (
	PlantainerTypeName = "plantainer"
	DeviceMQTTManager = &mcdevicemqttmanager.DeviceMQTTManagerSt{}
	DeviceRPCManager = mcdevicerpcmanager.CreateDeviceRPCManager("plantainerServerId", common.DevicesCollectionManager, DeviceMQTTManager)
)

func Init(dbsession *mgo.Session, dbName string) {

	InitMainModules(dbsession, dbName)
	InitRPCManager()
	InitMQTT()
}

func InitMQTT() {
	opts := mcmqttrouter.CreateConnOpts("tcp://localhost:1883", "randomString123qweasd", true)

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("OMG!!!! CONNECTION LOST BEACUSE: " + err.Error())
		DeviceMQTTManager.ReInitMQTT()
	}

	mqttClient := mcmqttrouter.CreateClient(opts)

	mqttRouter := mcmqttrouter.NewMQTTRouter(mqttClient, 1)
	mqttMainG := mqttRouter.Group("plantainerServerId")

	DeviceMQTTManager.Init(mqttMainG)

	DeviceMQTTManager.SetReqHandler(DeviceRPCManager.RPCReqHandler)

	fmt.Println("MQTT IS ACTIVATED")
}

func InitMainModules(dbsession *mgo.Session, dbName string) {
	common.InitDeviceColManager(dbsession, dbName)
}

func InitRPCManager() {
	DeviceRPCManager.AddDeviceCtrl(PlantainerTypeName, CreateNewPlantainerCtrl(PlantainerTypeName))
}

