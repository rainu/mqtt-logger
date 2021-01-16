package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

type MqttClient struct {
	BrokerAddress  string
	ClientId       string
	Username       string
	Password       string
	CaFile         string
	Timeout        time.Duration
	Topics         []string
	TopicBlacklist []*regexp.Regexp

	connectedOnce bool
	client        mqtt.Client
}

func (m *MqttClient) Connect() error {
	opts := mqtt.NewClientOptions()

	if m.CaFile != "" {
		certPool := x509.NewCertPool()
		certFile, err := ioutil.ReadFile(m.CaFile)
		if err != nil {
			return err
		}
		ok := certPool.AppendCertsFromPEM(certFile)
		if !ok {
			return fmt.Errorf("unable to parse ca cert '%s'", m.CaFile)
		}

		opts.SetTLSConfig(&tls.Config{
			RootCAs: certPool,
		})
	}

	opts.SetUsername(m.Username).
		SetPassword(m.Password).
		SetOnConnectHandler(m.handleConnect).
		SetConnectionLostHandler(m.handleDisconnect).
		SetAutoReconnect(true).
		AddBroker(m.BrokerAddress).
		SetClientID(m.ClientId)

	if m.Timeout > 0 {
		opts.SetConnectTimeout(m.Timeout)
	}

	m.client = mqtt.NewClient(opts)
	token := m.client.Connect()

	if token.Wait() {
		return nil
	}
	return token.Error()
}

func (m *MqttClient) handleConnect(_ mqtt.Client) {
	if m.connectedOnce {
		zap.L().Info("Reconnect to broker.")
	} else {
		zap.L().Info("Connection established.")
		m.connectedOnce = true
	}

	if err := m.subscribeTopics(); err != nil {
		zap.L().With(zap.Error(err)).Error("Error while subscribe topics!")
		os.Exit(2)
	}
}

func (m *MqttClient) handleDisconnect(_ mqtt.Client, err error) {
	zap.L().With(zap.Error(err)).Error("Connection to broker lost. Reconnecting...")
}

func (m *MqttClient) handleMessage(client mqtt.Client, message mqtt.Message) {
	if m.isBlacklisted(message.Topic()) {
		return
	}

	zap.L().With(
		zap.Bool("retained", message.Retained()),
		zap.Uint16("message_id", message.MessageID()),
		zap.String("Topics", message.Topic()),
		zap.ByteString("payload", message.Payload()),
	).Info("Incoming message.")
}

func (m *MqttClient) subscribeTopics() error {
	for _, topic := range m.Topics {
		token := m.client.Subscribe(topic, 2, m.handleMessage)
		if !token.WaitTimeout(5 * time.Second) {
			return fmt.Errorf("timeout while subscribe topic '%s': %w", topic, token.Error())
		}
	}

	return nil
}

func (m *MqttClient) isBlacklisted(topic string) bool {
	for _, regEx := range m.TopicBlacklist {
		if regEx.MatchString(topic) {
			return true
		}
	}

	return false
}
