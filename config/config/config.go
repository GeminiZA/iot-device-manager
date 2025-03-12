package config

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseType int

const (
	_ DatabaseType = iota
	POSTGRES
	SQLITE
)

type Config struct {
	SqlitePath       string
	ApiPort          string
	MQTTBrokerIp     string
	MQTTBrokerPort   string
	MQTTClientId     string
	MQTTUsername     string
	MQTTPassword     string
	MQTTUpdatesTopic string
	LogLevel         string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	path := os.Getenv("SQLITE_DATABASE_PATH")
	if path == "" {
		return nil, fmt.Errorf("missing SQLITE_DATABASE_PATH in env")
	}
	cfg.SqlitePath = path
	port := os.Getenv("API_PORT")
	if port == "" {
		return nil, fmt.Errorf("missing API_PORT in env")
	}
	cfg.ApiPort = port
	mqttBrokerPort := os.Getenv("MQTT_BROKER_PORT")
	if mqttBrokerPort == "" {
		return nil, fmt.Errorf("missing MQTT_BROKER_PORT in env")
	}
	portInt, err := strconv.Atoi(mqttBrokerPort)
	if err != nil || portInt < 1 || portInt > 50000 {
		return nil, fmt.Errorf("invalid MQTT_BROKER_PORT in env")
	}
	cfg.MQTTBrokerPort = mqttBrokerPort
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		return nil, fmt.Errorf("missing MQTT_BROKER_HOST in env")
	}
	ip := net.ParseIP(mqttBrokerHost)
	if ip == nil {
		ips, err := net.LookupIP(mqttBrokerHost)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve MQTT_BROKER_HOST ip, %v, %s", err, mqttBrokerHost)
		}
		if len(ips) > 0 {
			cfg.MQTTBrokerIp = ips[0].String()
		}
	} else {
		cfg.MQTTBrokerIp = ip.String()
	}
	mqttClientId := os.Getenv("MQTT_CLIENT_ID")
	if mqttClientId == "" {
		return nil, fmt.Errorf("missing MQTT_CLIENT_ID in env")
	}
	cfg.MQTTClientId = mqttClientId
	mqttUsername := os.Getenv("MQTT_USERNAME")
	if mqttUsername == "" {
		return nil, fmt.Errorf("missing MQTT_USERNAME in env")
	}
	cfg.MQTTUsername = mqttUsername
	mqttPassword := os.Getenv("MQTT_PASSWORD")
	if mqttPassword == "" {
		return nil, fmt.Errorf("missing MQTT_PASSWORD in env")
	}
	cfg.MQTTPassword = mqttPassword
	mqttUpdatesTopic := os.Getenv("MQTT_UPDATES_TOPIC")
	if mqttUpdatesTopic == "" {
		return nil, fmt.Errorf("missing MQTT_UPDATES_TOPIC in env")
	}
	cfg.MQTTUpdatesTopic = mqttUpdatesTopic
	loglevel := os.Getenv("LOG_LEVEL")
	if loglevel == "" {
		cfg.LogLevel = "ERROR"
	} else {
		cfg.LogLevel = loglevel
	}
	return &cfg, nil
}
