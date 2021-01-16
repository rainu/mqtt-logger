package main

import (
	"github.com/rainu/mqtt-logger/internal/mqtt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// read config from ENV
	c := NewConfig()
	c.Log()

	client := mqtt.MqttClient{
		BrokerAddress:  c.MqttBrokerAddress,
		ClientId:       "mqtt-logger",
		Username:       c.MqttUsername,
		Password:       c.MqttPassword,
		CaFile:         c.MqttCaCert,
		Timeout:        c.MqttTimeout,
		Topics:         c.MqttTopics,
		TopicBlacklist: c.MqttTopicBlacklist,
	}
	if err := client.Connect(); err != nil {
		zap.L().With(zap.Error(err)).Error("Error while connecting to mqtt broker!")
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	//wait until interrupt
	<-stop
}
