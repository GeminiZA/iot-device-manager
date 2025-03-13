package routes

import (
	"strconv"

	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/gofiber/fiber/v2"
)

func (handler *HttpHandler) GetAllDevices(c *fiber.Ctx) error {
	params := c.AllParams()
	deviceIdStr, ok := params["id"]
	if !ok {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing device id from parameters",
		})

	}
	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id; must be a number",
		})

	}
	device, err := handler.dr.GetDevice(uint(deviceId))
	if device == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
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
