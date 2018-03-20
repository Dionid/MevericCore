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
	"mevericcore/mccommon"
)

var (
	PlantainerTypeName = "plantainer"
	DeviceMQTTManager = &mcdevicemqttmanager.DeviceMQTTManagerSt{}
	DeviceRPCManager = mcdevicerpcmanager.CreateDeviceRPCManager("plantainerServerId", DevicesCollectionManager, DeviceMQTTManager)
)

func Init(userColMan *mccommon.UsersCollectionManagerSt, dbsession *mgo.Session, dbName string, e *echo.Group) {

	InitMainModules(dbsession, dbName)
	InitRPCManager()
	InitMQTT()

	InitModelHttp(e)
}

func InitMQTT() {
	//opts := mcmqttrouter.CreateConnOpts("tcp://iot.eclipse.org:1883", "randomString123qweasd", true)
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
	initDeviceColManager(dbsession, dbName)
}

func InitRPCManager() {
	DeviceRPCManager.AddDeviceCtrl(PlantainerTypeName, CreateNewPlantainerCtrl(PlantainerTypeName))
}

func InitModelHttp(e *echo.Group) {
	UserPlantainerController := &UserPlantainerControllerSt{}
	mcecho.CreateModelControllerRoutes(e, "/plantainer", UserPlantainerController)
}