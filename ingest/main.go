package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	ConnectMQTT()
	SubscribeECU()
	SubscribeBattery()
	ConnectDB()
	StartServer()
}

var Client mqtt.Client

var MQTTHost = "localhost"
var MQTTPort = "1883"

func ConnectMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", MQTTHost, MQTTPort))
	opts.OnConnect = connectHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	Client = client
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("[MQ] Connected to MQTT broker!")
}

func SubscribeECU() {
	Client.Subscribe("ingest/ecu", 0, func(client mqtt.Client, msg mqtt.Message) {
		ecu := ECUFromBytes(msg.Payload())
		ecu, err := CreateECU(ecu)
		if err != nil {
			fmt.Printf("failed to create ecu: %v", err)
		}
		fmt.Printf("[MQ] Received ecu message: %v\n", ecu)
	})
}

func SubscribeBattery() {
	Client.Subscribe("ingest/battery", 0, func(client mqtt.Client, msg mqtt.Message) {
		battery := BatteryFromBytes(msg.Payload())
		battery, err := CreateBattery(battery)
		if err != nil {
			fmt.Printf("failed to create battery: %v", err)
		}
		fmt.Printf("[MQ] Received battery message: %v\n", battery)
	})
}
