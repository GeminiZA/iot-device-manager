package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GeminiZA/iot-device-manager/config/config"
	"github.com/GeminiZA/iot-device-manager/config/database"
	"github.com/GeminiZA/iot-device-manager/controllers/api"
	"github.com/GeminiZA/iot-device-manager/controllers/mqtt"
	"github.com/GeminiZA/iot-device-manager/controllers/timer"
	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/GeminiZA/iot-device-manager/view/mqttHandlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Building with sqlite so long, implementing postgres should be easy
	sqliteDb, err := database.ConnectSqlite(cfg.SqlitePath)
	if err != nil {
		log.Fatal(err)
	}
	dr := models.NewDeviceRepository(sqliteDb)

	// Start the timeout worker
	timer := timer.NewTimer(dr)
	err = timer.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer timer.Stop()

	// Start mqtt Client
	mqttClient, err := mqtt.Connect(mqtt.MQTTConfig{
		Broker:           fmt.Sprintf("%s:%s", cfg.MQTTBrokerIp, cfg.MQTTBrokerPort),
		ClientId:         cfg.MQTTClientId,
		Username:         cfg.MQTTUsername,
		Password:         cfg.MQTTPassword,
		UpdatesTopicPath: cfg.MQTTUpdatesTopic,
		DeviceRepository: dr,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer mqttClient.Disconnect()

	mqttHandler := mqttHandlers.NewHandler(mqttClient, dr, cfg.MQTTUpdatesTopic)
	err = mqttHandler.SubscribeDeviceTelemetry()
	if err != nil {
		log.Fatal(err)
	}

	// Start http api
	apiCfg, err := api.NewApiConfig(cfg.ApiPort, dr, mqttHandler, timer)
	if err != nil {
		log.Fatal(err)
	}
	apiCfg.Listen()

	//Gracefully exit

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	apiCfg.Stop()
	mqttClient.Disconnect()
	fmt.Println("Successfully shut down")
	os.Exit(0)
}
