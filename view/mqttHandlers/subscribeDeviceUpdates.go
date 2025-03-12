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
	var device struct {
		Status    string         `json:"status"`
		Telemetry datatypes.JSON `json:"telemetry"`
	}
	err = json.Unmarshal(msg.Payload(), &device)
	if err != nil {
		log.Errorf("error handling mqtt update message: invalid message payload: %v", err)
		return
	}
	newDeviceDetails := models.NewDeviceDetails{
		Status:    device.Status,
		Telemetry: device.Telemetry,
	}
	err = handler.dr.UpdateDevice(uint(deviceID), &newDeviceDetails)
	if err != nil {
		log.Errorf("error handling mqtt update message: cannot update database: %v", err)
		return
	}
	log.Debugf("Got update from mqtt %v", newDeviceDetails)
}

func (handler *MQTTHandler) SubscribeDeviceUpdates() error {
	topic := path.Join(handler.updatesTopicPath, "+")
	fmt.Println(topic)
	return handler.client.Subscribe(topic, handler.handleUpdateMessage)
}
