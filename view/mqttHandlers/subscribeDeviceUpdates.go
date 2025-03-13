package mqttHandlers

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/GeminiZA/iot-device-manager/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/datatypes"
)

func (handler *MQTTHandler) handleUpdateMessage(c mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	topicComponenets := strings.Split(topic, "/")
	if len(topicComponenets) < 2 {
		log.Errorf("error handling mqtt update message: invalid topic: %s", topic)
		return
	}
	deviceIDStr := topicComponenets[len(topicComponenets)-1]
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		log.Errorf("error handling mqtt update message: invalid device id: %s", deviceIDStr)
		return
	}
	update := struct {
		Status    string         `json:"status"`
		Telemetry datatypes.JSON `json:"telemetry"`
	}{}
	err = json.Unmarshal(msg.Payload(), &update)
	if err != nil {
		log.Errorf("error handling mqtt update message: invalid message payload: %v", err)
		return
	}
	err = handler.dr.AddTelemetry(uint(deviceID), &update.Telemetry)
	if err != nil {
		log.Errorf("error handling mqtt update message: cannot update database: %v", err)
		return
	}
	err = handler.dr.UpdateDevice(uint(deviceID), &models.NewDeviceDetails{
		Status: update.Status,
	})
	if err != nil {
		log.Errorf("error handling mqtt update message: cannot update database: %v", err)
		return
	}
	log.Debugf("Got update from mqtt %v", update)
}

func (handler *MQTTHandler) SubscribeDeviceTelemetry() error {
	topic := path.Join(handler.updatesTopicPath, "+")
	fmt.Println(topic)
	return handler.client.Subscribe(topic, handler.handleUpdateMessage)
}
