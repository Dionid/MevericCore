package mcdevicemqttmanager

import (
	"mevericcore/mcmqttrouter"
)

//var (
//	DeviceMQTTManager = &DeviceMQTTManagerSt{}
//)

func (this *DeviceMQTTManagerSt) Init(router *mcmqttrouter.MQTTRouter) {
	this.InitMQTT(router)
}

func (this *DeviceMQTTManagerSt) InitMQTT(router *mcmqttrouter.MQTTRouter) {
	this.router = router
	this.SubscribeMainRoutes()
}

func (this *DeviceMQTTManagerSt) ReInitMQTT() {
	this.UnSubscribeFromAll()
	this.SubscribeMainRoutes()
}