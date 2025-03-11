package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConfig struct {
	Broker         string
	ClientId       string
	Username       string
	Password       string
	ConnectTimeout time.Duration
}

type MQTTClient struct {
	mqttClient mqtt.Client
}

func Connect(cfg *MQTTConfig) (*MQTTClient, error) {
	// Set timeout if not specified
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 3 * time.Second
	}

	clientOpts, err := createClientOpts(cfg)
	if err != nil {
		return nil, err
	}
	mqttClient := mqtt.NewClient(clientOpts)
	token := mqttClient.Connect()
	for !token.WaitTimeout(cfg.ConnectTimeout) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}
	client := MQTTClient{
		mqttClient: mqttClient,
	}
	return &client, nil
}

func (client *MQTTClient) Disconnect() {
	client.mqttClient.Disconnect(250)
}

func createClientOpts(cfg *MQTTConfig) (*mqtt.ClientOptions, error) {
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
	return opts, nil
}
