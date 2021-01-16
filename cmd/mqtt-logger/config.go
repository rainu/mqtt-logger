package main

import (
	"github.com/alexflint/go-arg"
	"go.uber.org/zap"
	"os"
	"regexp"
	"time"
)

type Config struct {
	MqttBrokerAddress     string        `arg:"--broker,required,env:MQTT_BROKER_ADDRESS,help:The mqtt broker address."`
	MqttTopics            []string      `arg:"--topic,separate,env:MQTT_TOPICS,help:List of the mqtt topics to subscribe to."`
	MqttRawTopicBlacklist []string      `arg:"--blacklist,env:MQTT_TOPIC_BLACKLIST,help:List of regular expression."`
	MqttUsername          string        `arg:"--user,env:MQTT_USERNAME,help:The mqtt username."`
	MqttPassword          string        `arg:"--password,env:MQTT_PASSWORD,help:The mqtt password."`
	MqttTimeout           time.Duration `default:"1m" arg:"--timeout,env:MQTT_TIMEOUT,help:The timeout for the mqtt connection."`
	MqttCaCert            string        `arg:"--ca-cert,env:MQTT_CA_CERT,help:The path of the CA-Cert file for secure mqtt connection."`

	MqttTopicBlacklist []*regexp.Regexp `arg:"-"`
}

func NewConfig() *Config {
	c := &Config{}
	arg.MustParse(c)

	if len(c.MqttTopics) == 0 {
		c.MqttTopics = []string{"#"}
	}

	for _, regExp := range c.MqttRawTopicBlacklist {
		compiled, err := regexp.Compile(regExp)
		if err != nil {
			zap.L().With(
				zap.Error(err),
				zap.String("expression", regExp),
			).Error("Error while compoiling blacklist regular expression")
			os.Exit(1)
		}

		c.MqttTopicBlacklist = append(c.MqttTopicBlacklist, compiled)
	}

	return c
}

func (s *Config) Log() {
	zap.L().With(
		zap.Strings("MQTT_TOPICS", s.MqttTopics),
		zap.Strings("MQTT_TOPIC_BLACKLIST", s.MqttRawTopicBlacklist),
		zap.String("MQTT_BROKER_ADDRESS", s.MqttBrokerAddress),
		zap.String("MQTT_USERNAME", s.MqttUsername),
		zap.String("MQTT_PASSWORD", "*hidden*"),
		zap.Duration("MQTT_TIMEOUT", s.MqttTimeout),
		zap.String("MQTT_CA_CERT", s.MqttCaCert),
	).Info("Using configuration")
}
