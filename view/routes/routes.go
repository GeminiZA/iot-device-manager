package routes

import (
	"fmt"

	"github.com/GeminiZA/iot-device-manager/controllers/timer"
	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/GeminiZA/iot-device-manager/view/mqttHandlers"
)

type HttpHandler struct {
	dr          *models.DeviceRepository
	mqttHandler *mqttHandlers.MQTTHandler
	timer       *timer.Timer
}

func NewHandler(dr *models.DeviceRepository, mqttHandler *mqttHandlers.MQTTHandler, timer *timer.Timer) (*HttpHandler, error) {
	if dr == nil {
		return nil, fmt.Errorf("missing device repository")
	}
	if mqttHandler == nil {
		return nil, fmt.Errorf("missing mqtt handler")
	}
	return &HttpHandler{
		dr:          dr,
		mqttHandler: mqttHandler,
		timer:       timer,
	}, nil
}
