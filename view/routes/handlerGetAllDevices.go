package routes

import (
	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/gofiber/fiber/v2"
)

func (handler *HttpHandler) GetAllDevices(c *fiber.Ctx) error {
	devices, err := handler.dr.GetAllDevices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	response := struct {
		Devices []models.Device
	}{
		Devices: *devices,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
