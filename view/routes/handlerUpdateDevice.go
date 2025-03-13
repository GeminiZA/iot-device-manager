package routes

import (
	"errors"
	"strconv"

	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (handler *HttpHandler) UpdateDevice(c *fiber.Ctx) error {
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
	c.Accepts("application/json")
	var request struct {
		Status string `json:"status"`
		Name   string `json:"name"`
	}
	err = c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	err = handler.dr.UpdateDevice(uint(deviceId), &models.NewDeviceDetails{
		Status: request.Status,
		Name:   request.Name,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = handler.mqttHandler.PublishDeviceUpdate(&models.Device{
		ID:     uint(deviceId),
		Status: request.Status,
		Name:   request.Name,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
