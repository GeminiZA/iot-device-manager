package routes

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func (handler *HttpHandler) GetDevice(c *fiber.Ctx) error {
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
	response := struct {
		ID        uint           `json:"id"`
		Name      string         `json:"name"`
		Status    string         `json:"status"`
		Telemetry datatypes.JSON `json:"telemetry"`
	}{
		ID:        device.ID,
		Name:      device.Name,
		Status:    device.Status,
		Telemetry: device.Telemetry,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
