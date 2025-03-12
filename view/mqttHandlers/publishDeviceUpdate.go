package mqttHandlers

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/GeminiZA/iot-device-manager/models"
	"gorm.io/datatypes"
)

func (handler *MQTTHandler) PublishDeviceUpdate(device *models.Device) error {
	retDevice := struct {
		Status    string         `json:"status"`
		Telemetry datatypes.JSON `json:"telemetry"`
	}{
		Status:    device.Status,
		Telemetry: device.Telemetry,
	}
	topic := path.Join(handler.updatesTopicPath, fmt.Sprintf("%d", device.ID))
	message, err := json.Marshal(retDevice)
	if err != nil {
		return err
	}
	return handler.client.Publish(topic, message)
}
