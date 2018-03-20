package mcmqttrouter

import (
	"github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
	"log"
	"fmt"
)

func CreateConnOpts(brokerName string, clientId string, debug bool) *mqtt.ClientOptions {
	if debug {
		//mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
	}
	opts := mqtt.NewClientOptions().AddBroker(brokerName).SetClientID(clientId)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetDefaultPublishHandler(standardHandler)
	opts.SetPingTimeout(30 * time.Second)
	opts.AutoReconnect = true
	opts.OnConnectionLost = func (c mqtt.Client, err error) {
		fmt.Println("OMG!!!! CONNECTION LOST BEACUSE: " + err.Error())
	}

	return opts
}

func CreateClient(opts *mqtt.ClientOptions) *mqtt.Client {
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return &c
}