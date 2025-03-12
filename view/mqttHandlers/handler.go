package mqttHandlers

import (
	"github.com/GeminiZA/iot-device-manager/controllers/mqtt"
	"github.com/GeminiZA/iot-device-manager/models"
)

type MQTTHandler struct {
	client           *mqtt.MQTTClient
	dr               *models.DeviceRepository
	updatesTopicPath string
}

func NewHandler(client *mqtt.MQTTClient, dr *models.DeviceRepository, updatesTopicPath string) *MQTTHandler {
	return &MQTTHandler{
		client:           client,
		dr:               dr,
		updatesTopicPath: updatesTopicPath,
	}
}
