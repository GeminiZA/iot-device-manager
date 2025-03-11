package main

import (
	"fmt"
	"log"

	"github.com/GeminiZA/iot-device-manager/config/config"
	"github.com/GeminiZA/iot-device-manager/config/database"
	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/GeminiZA/iot-device-manager/view/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	var db *gorm.DB
	switch cfg.DatabaseType {
	case config.POSTGRES:
		// To implement
	case config.SQLITE:
		// Building with sqlite so long, implementing postgres should be easy
		sqliteDb, err := database.ConnectSqlite(cfg.SqlitePath)
		if err != nil {
			log.Fatal(err)
		}
		db = sqliteDb
	}
	dr := models.NewDeviceRepository(db)
	app := fiber.New()
	handler := routes.Handler{
		Dr: dr,
	}
	app.Put("/assets/:id", handler.UpdateDevice)
	app.Get("/assets/:id", handler.GetDevice)
	app.Post("/assets", handler.UniqueDeviceIDMiddleware(), handler.CreateDevice)
	app.Delete("/assets/:id", handler.DeleteDevice)
	app.Listen(fmt.Sprintf(":%s", cfg.ApiPort))
}
