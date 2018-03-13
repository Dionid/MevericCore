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

func InitMQTT(router *mcmqttrouter.MQTTRouter) {
	DeviceMQTTManager.router = router
	DeviceMQTTManager.SubscribeMainRoutes()
}

func ReInitMQTT() {
	DeviceMQTTManager.UnSubscribeFromAll()
	DeviceMQTTManager.SubscribeMainRoutes()
}