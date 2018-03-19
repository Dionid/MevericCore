package mcplantainer

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"mevericcore/mcmqttrouter"
	"mevericcore/mcdevicemqttmanager"
	"gopkg.in/mgo.v2"
	"mevericcore/mcdevicerpcmanager"
	"mevericcore/mcecho"
	"github.com/labstack/echo"
)

var (
	PlantainerTypeName = "smartHeater"
	DeviceMQTTManager = mcdevicemqttmanager.DeviceMQTTManager
	DeviceRPCManager = mcdevicerpcmanager.CreateDeviceRPCManager("plantainerServerId", DevicesCollectionManager, DeviceMQTTManager)
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
	DeviceMQTTManager.SetReqHandler(DeviceRPCManager.RPCReqHandler)

	fmt.Println("MQTT IS ACTIVATED")
}

func Init(dbsession *mgo.Session, dbName string) {
	initDeviceColManager(dbsession, dbName)

	DeviceRPCManager.AddDeviceCtrl(PlantainerTypeName, CreateNewPlantainerCtrl(PlantainerTypeName))

	activateMQTT()
}

func InitHttp(dbsession *mgo.Session, dbName string, e *echo.Group) {
	initDeviceColManager(dbsession, dbName)

	UserPlantainerController := &UserPlantainerControllerSt{}

	mcecho.CreateModelControllerRoutes(e, "/plantainer", UserPlantainerController)
}