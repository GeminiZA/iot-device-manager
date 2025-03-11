package routes

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (handler *Handler) UpdateDevice(c *fiber.Ctx) error {
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
		Status    string         `json:"status"`
		Telemetry datatypes.JSON `json:"telemetry"`
	}
	err = c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	var telemetry datatypes.JSON
	err = json.Unmarshal([]byte(request.Telemetry), &telemetry)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = handler.Dr.UpdateDevice(uint(deviceId), &models.NewDeviceDetails{
		Telemetry: telemetry,
		Status:    request.Status,
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
	return c.SendStatus(fiber.StatusOK)
}
