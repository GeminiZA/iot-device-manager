package routes

import (
	"github.com/gofiber/fiber/v2"
)

func (handler *HttpHandler) UniqueDeviceIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request struct {
			DeviceID uint `json:"id"`
		}
		err := c.BodyParser(&request)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		_, err = handler.dr.GetDevice(request.DeviceID)
		if err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "id must be unique",
			})
		}
		return c.Next()
	}
}
