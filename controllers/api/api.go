package api

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/GeminiZA/iot-device-manager/controllers/timer"
	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/GeminiZA/iot-device-manager/view/mqttHandlers"
	"github.com/GeminiZA/iot-device-manager/view/routes"
	"github.com/gofiber/fiber/v2"
)

type APIConfig struct {
	port             string
	deviceRepository *models.DeviceRepository
	fiberApp         *fiber.App
	timer            *timer.Timer
}

func NewApiConfig(port string, deviceRepository *models.DeviceRepository, mqttHandler *mqttHandlers.MQTTHandler, timer *timer.Timer) (*APIConfig, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	if portInt < 0 || portInt > 49151 {
		return nil, fmt.Errorf("invalid port")
	}
	cfg := APIConfig{
		port:             port,
		deviceRepository: deviceRepository,
		timer:            timer,
	}
	cfg.fiberApp = fiber.New(fiber.Config{
		AppName:     "iot-device-manager",
		Concurrency: 256 * 1024, // Default
	})
	handler, err := routes.NewHandler(deviceRepository, mqttHandler, timer)
	if err != nil {
		return nil, err
	}
	cfg.fiberApp.Put("/assets/:id", handler.UpdateDevice)
	cfg.fiberApp.Get("/assets/:id", handler.GetDevice)
	cfg.fiberApp.Get("/assets", handler.GetAllDevices)
	cfg.fiberApp.Post("/assets", handler.CreateDevice)
	cfg.fiberApp.Delete("/assets/:id", handler.DeleteDevice)
	return &cfg, nil
}

func (cfg *APIConfig) run() {
	log.Fatal(cfg.fiberApp.Listen(fmt.Sprintf(":%s", cfg.port)))
}

func (cfg *APIConfig) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cfg.fiberApp.ShutdownWithContext(ctx); err != nil {
		return err
	}
	return nil
}

func (cfg *APIConfig) Listen() error {
	if cfg.fiberApp == nil {
		return fmt.Errorf("app not initialized, call NewApiConfig first")
	}
	go cfg.run()
	return nil
}
