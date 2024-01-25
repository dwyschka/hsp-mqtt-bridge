package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"strconv"
	"time"
)

func initMqtt() mqtt.Client {

	clientId := "hsp_pellet_stove"
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Println("Connected")
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Printf("Connect lost: %v", err)

	}

	options := mqtt.NewClientOptions()
	broker := os.Getenv("MQTT_IP")
	port, _ := strconv.Atoi(os.Getenv("MQTT_PORT"))
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	options.SetClientID(clientId)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectLostHandler

	useAuth, _ := strconv.ParseBool(os.Getenv("MQTT_USE_AUTH"))
	if useAuth {
		options.Username = os.Getenv("MQTT_USER")
		options.Password = os.Getenv("MQTT_PASSWORD")
	}
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Got Error: %v\r\n", token.Error())
		log.Println("Try again in 5 seconds..")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
	return client
}
