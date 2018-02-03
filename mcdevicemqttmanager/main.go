package mcdevicemqttmanager

import (
	"mevericcore/mcmqttrouter"
)

var (
	DeviceMQTTManager = &DeviceMQTTManagerSt{}
)

func Init(router *mcmqttrouter.MQTTRouter) {
	InitMQTT(router)
}

//func main() {
//	// 1. Create MQTT broker connection
//	// 2. DeviceMQTTManagerSt = &DeviceMQTTManagerSt{}
//	// 3. InitMQTT(router)
//	// 4. Subscribe to QueueManager
//}
