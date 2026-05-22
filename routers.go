package main

import (
	"go-cloud/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func startRoutingApp(app *fiber.App, handlerStruct *handler.Handler) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://127.0.0.1:5173",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint: "/health",
	}))

	app.Get("/monitor", monitor.New())
	app.Get("/checkdb", handlerStruct.CheckDB)
	app.Post("/users", handlerStruct.AddNewUser)
}
