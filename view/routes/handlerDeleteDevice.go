package routes

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (handler *HttpHandler) DeleteDevice(c *fiber.Ctx) error {
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
	err = handler.dr.DeleteDevice(uint(deviceId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = handler.mqttHandler.PublishDeviceCommand(uint(deviceId), "stop")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
