package device

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"mevericcore/mcmqttrouter"
	"mevericcore/mcdevicemqttmanager"
	"gopkg.in/mgo.v2"
	"mevericcore/mcdevicerpcmanager_old"
	"mevericcore/mcplantainer/common"
	"mevericcore/mccommon"
	"mevericcore/mcinnerrpc"
)

var (
	PlantainerServerId = "plantainerServerId"

	InnerRPCMan                                              = mcinnerrpc.New()
	DeviceMQTTMan                                            = &mcdevicemqttmanager.DeviceMQTTManagerSt{}
	DeviceRPCMan  *mcdevicerpcmanager_old.DeviceRPCManagerSt = nil
)

func Init(dbsession *mgo.Session, dbName string) {
	InitInnerRPCManager()

	InitMainModules(dbsession, dbName)

	InitRPCManager()

	InitMQTT()
}

func InitInnerRPCManager() {
	InnerRPCMan.Init()
	InnerRPCMan.Service.Subscribe("Device.RPC.Send", func(msg *mcinnerrpc.Msg) {
		rpcData := mccommon.RPCMsg{}

		if err := rpcData.UnmarshalJSON(msg.Data); err != nil {
			return
		}

		DeviceMQTTMan.PublishJSON(rpcData.Dst + "/rpc", rpcData)
	})
}

func InitMainModules(dbsession *mgo.Session, dbName string) {
	common.InitDeviceColManager(dbsession, dbName)
}

func InitRPCManager() {
	DeviceRPCMan = mcdevicerpcmanager_old.CreateDeviceRPCManager(PlantainerServerId, common.PlantainerCollectionManager, InnerRPCMan.SendRPCMsgToUser)
	DeviceRPCMan.AddDeviceCtrl(CreateNewPlantainerCtrl())
}

func InitMQTT() {
	opts := mcmqttrouter.CreateConnOpts("tcp://localhost:1883", "randomString123qweasd", true)

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("OMG!!!! CONNECTION LOST BEACUSE: " + err.Error())
		DeviceMQTTMan.ReInitMQTT()
	}

	mqttClient := mcmqttrouter.CreateClient(opts)

	mqttRouter := mcmqttrouter.NewMQTTRouter(mqttClient, 1)
	mqttMainG := mqttRouter.Group(PlantainerServerId)

	DeviceMQTTMan.Init(mqttMainG)
	DeviceMQTTMan.SetReqHandler(DeviceRPCMan.RPCReqHandler)
}

