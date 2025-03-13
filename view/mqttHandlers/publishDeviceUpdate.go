package mqttHandlers

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/GeminiZA/iot-device-manager/models"
)

func (handler *MQTTHandler) PublishDeviceUpdate(device *models.Device) error {
	retDevice := struct {
		Status string `json:"status"`
		Name   string `json:"name"`
	}{
		Status: device.Status,
		Name:   device.Name,
	}
	topic := path.Join(handler.updatesTopicPath, fmt.Sprintf("%d", device.ID))
	message, err := json.Marshal(retDevice)
	if err != nil {
		return err
	}
	return handler.client.Publish(topic, message)
}

func (handler *MQTTHandler) PublishDeviceCommand(deviceId uint, command string) error {
	retMessage := struct {
		Command string `json:"command"`
	}{
		Command: command,
	}
	topic := path.Join(handler.updatesTopicPath, fmt.Sprintf("%d", deviceId))
	message, err := json.Marshal(retMessage)
	if err != nil {
		return err
	}
	return handler.client.Publish(topic, message)
}
