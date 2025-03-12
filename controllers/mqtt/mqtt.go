package mqtt

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"time"

	"github.com/GeminiZA/iot-device-manager/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2/log"
)

type MQTTConfig struct {
	Broker           string
	ClientId         string
	Username         string
	Password         string
	ConnectTimeout   time.Duration
	UpdatesTopicPath string
	PublishTimeout   time.Duration
	DeviceRepository *models.DeviceRepository
}

type MQTTClient struct {
	Connected  bool
	mqttClient mqtt.Client
	cfg        MQTTConfig
}

func Connect(cfg MQTTConfig) (*MQTTClient, error) {
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 3 * time.Second
	}
	if cfg.PublishTimeout == 0 {
		cfg.PublishTimeout = 500 * time.Millisecond
	}
	if cfg.UpdatesTopicPath == "" {
		return nil, fmt.Errorf("missing updates topic path")
	}
	if cfg.DeviceRepository == nil {
		return nil, fmt.Errorf("missing DeviceRepository")
	}
	client := MQTTClient{
		cfg:       cfg,
		Connected: false,
	}
	// Set options
	if cfg.Broker == "" ||
		cfg.ClientId == "" ||
		cfg.Username == "" ||
		cfg.Password == "" {
		return nil, fmt.Errorf("missing config settings")
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientId)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(2 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Errorf("MQTT Broker connection lost (%v)...", err)
		client.Connected = false
	})
	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		log.Debugf("Connecting to MQTT broker(%s)...", broker.String())
		return tlsCfg
	})
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Debugf("Successfully connected to MQTT broker")
		client.Connected = true
	})

	mqttClient := mqtt.NewClient(opts)
	token := mqttClient.Connect()
	for !token.WaitTimeout(cfg.ConnectTimeout) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	client.mqttClient = mqttClient
	return &client, nil
}

func (client *MQTTClient) Publish(topic string, message []byte) error {
	token := client.mqttClient.Publish(topic, 1, false, message)
	for !token.WaitTimeout(client.cfg.PublishTimeout) {
	}
	if err := token.Error(); err != nil {
		return err
	}
	return nil
}

func (client *MQTTClient) Subscribe(topic string, callback mqtt.MessageHandler) error {
	token := client.mqttClient.Subscribe(topic, 1, callback)
	for !token.WaitTimeout(client.cfg.PublishTimeout) {
	}
	if err := token.Error(); err != nil {
		return err
	}
	return nil
}

func (client *MQTTClient) Disconnect() {
	client.mqttClient.Disconnect(250)
}
