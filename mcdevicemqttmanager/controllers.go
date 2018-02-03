package mcdevicemqttmanager

import (
	"mevericcore/mcmqttrouter"
)

func InitMQTT(router *mcmqttrouter.MQTTRouter) {
	DeviceMQTTManager.router = router
	DeviceMQTTManager.SubscribeMainRoutes()
}

func ReInitMQTT() {
	DeviceMQTTManager.UnSubscribeFromAll()
	DeviceMQTTManager.SubscribeMainRoutes()
}
