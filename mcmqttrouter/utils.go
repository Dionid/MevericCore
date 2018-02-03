package mcmqttrouter

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"os"
)

func unsubscribeFromChannels(c *mqtt.Client, topics ...string) {
	if token := (*c).Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	fmt.Println("Usubscribed from: ")
	for _, t := range topics {
		fmt.Println("Topic name: " + t)
	}
}

func addSubscribeRoute(c *mqtt.Client, topicName string, handler mqtt.MessageHandler)	 {
	if token := (*c).Subscribe(topicName, 0, handler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func standardHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}