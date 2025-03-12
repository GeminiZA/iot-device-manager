package routes

import (
	"fmt"

	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/GeminiZA/iot-device-manager/view/mqttHandlers"
)

type HttpHandler struct {
	dr          *models.DeviceRepository
	mqttHandler *mqttHandlers.MQTTHandler
}

func NewHandler(dr *models.DeviceRepository, mqttHandler *mqttHandlers.MQTTHandler) (*HttpHandler, error) {
	if dr == nil {
		return nil, fmt.Errorf("missing device repository")
	}
	if mqttHandler == nil {
		return nil, fmt.Errorf("missing mqtt handler")
	}
	return &HttpHandler{
		dr:          dr,
		mqttHandler: mqttHandler,
	}, nil
}
