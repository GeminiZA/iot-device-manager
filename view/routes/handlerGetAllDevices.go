package routes

import (
	"github.com/gofiber/fiber/v2"
)

type RetDevice struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (handler *HttpHandler) GetAllDevices(c *fiber.Ctx) error {
	devices, err := handler.dr.GetAllDevices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	retDevices := make([]RetDevice, 0)
	for _, device := range devices {
		retDevices = append(retDevices, RetDevice{
			Id:     device.ID,
			Name:   device.Name,
			Status: device.Status,
		})
	}

	response := struct {
		Devices []RetDevice `json:"devices"`
	}{
		Devices: retDevices,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
