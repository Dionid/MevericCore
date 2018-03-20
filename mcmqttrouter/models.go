package mcmqttrouter

import "github.com/eclipse/paho.mqtt.golang"

type MQTTRouter struct {
	Client *mqtt.Client
	prefix string
	Routes []string
	standardQOS byte
}

func NewMQTTRouter(c *mqtt.Client, standardQOS byte) *MQTTRouter {
	return &MQTTRouter{
		c,
		"",
		[]string{},
		standardQOS,
	}
}

func (this *MQTTRouter) Subscribe(path string, handler mqtt.MessageHandler) {
	this.Routes = append(this.Routes, path)
	addSubscribeRoute(this.Client, this.prefix + path, handler)
}

func (this *MQTTRouter) UnSubscribeFromAll() {
	unsubscribeFromChannels(this.Client, this.Routes...)
}

func (this *MQTTRouter) PublishCustom(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return (*this.Client).Publish(topic, qos, retained, payload)
}

func (this *MQTTRouter) Publish(topic string, payload interface{}) mqtt.Token {
	return (*this.Client).Publish(topic, this.standardQOS, false, payload)
}

func (this *MQTTRouter) Group(path string) *MQTTRouter{
	return &MQTTRouter{this.Client, this.prefix + path, []string{}, this.standardQOS}
}
