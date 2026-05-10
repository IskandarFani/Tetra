package handler

import (
	"go-cloud/database"

	"github.com/gofiber/contrib/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func Health() fiber.Handler {
	return healthcheck.New(healthcheck.Config{
		ResponseFormat: healthcheck.FormatJSON,
	})
}

func Monitor() fiber.Handler {
	return monitor.New()
}

func CheckDB() fiber.Handler {

	var status database.DBStatus

	status = database.Connect()

	if !status.Connected {
		return func(c fiber.Ctx) error {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": status.Message,
			})
		}
	}
	return func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": status.Message,
		})
	}
}
