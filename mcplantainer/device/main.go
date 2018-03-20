package device

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"mevericcore/mcmqttrouter"
	"mevericcore/mcdevicemqttmanager"
	"gopkg.in/mgo.v2"
	"mevericcore/mcdevicerpcmanager"
	"mevericcore/mcplantainer/common"
	"github.com/nats-io/go-nats"
	"mevericcore/mccommon"
)

var (
	PlantainerTypeName = "plantainer"
	NATSCon *nats.Conn = nil
	DeviceMQTTManager = &mcdevicemqttmanager.DeviceMQTTManagerSt{}
	DeviceRPCManager *mcdevicerpcmanager.DeviceRPCManagerSt = nil
)

func SendToUser(msg *mccommon.RPCMsg) error {
	if bData, err := msg.MarshalJSON(); err != nil {
		return err
	} else {
		NATSCon.Publish("User.RPC.Send", bData)
	}

	return nil
}

func Init(dbsession *mgo.Session, dbName string) {
	nc, _ := nats.Connect(nats.DefaultURL)
	NATSCon = nc

	NATSCon.Subscribe("Device.RPC.Send", func(msg *nats.Msg) {
		rpcData := mccommon.RPCMsg{}

		if err := rpcData.UnmarshalJSON(msg.Data); err != nil {
			return
		}

		DeviceMQTTManager.PublishJSON(rpcData.Dst + "/rpc", rpcData)
	})

	DeviceRPCManager = mcdevicerpcmanager.CreateDeviceRPCManager("plantainerServerId", common.PlantainerCollectionManager, DeviceMQTTManager, SendToUser)

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

