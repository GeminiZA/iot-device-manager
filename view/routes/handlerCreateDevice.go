package routes

import (
	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/gofiber/fiber/v2"
)

func (handler *HttpHandler) CreateDevice(c *fiber.Ctx) error {
	var request struct {
		DeviceName string `json:"name"`
		DeviceId   uint   `json:"id"`
	}
	err := c.BodyParser(&request)
	if err != nil || request.DeviceName == "" || request.DeviceId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	newDevice := models.NewDevice(request.DeviceName, request.DeviceId)
	err = handler.dr.CreateDevice(newDevice)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Return device name and id
	response := struct {
		Name string `json:"name"`
		ID   uint   `json:"id"`
	}{
		Name: newDevice.Name,
		ID:   newDevice.ID,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}
