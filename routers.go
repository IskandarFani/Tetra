package main

import (
	"go-cloud/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func startRoutingApp(app *fiber.App, h *handler.Handler) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://127.0.0.1:5173",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint: "/health",
	}))

	app.Get("/monitor", monitor.New())
	app.Get("/checkdb", h.CheckDB)

	auth := app.Group("/auth")

	auth.Post("/register", h.Register)
	auth.Post("/login", h.LogIn)
	auth.Post("/refresh", h.ReadRefreshToken, h.RefreshToken)

	api := app.Group("/api")

	api.Get("/profile", h.ReadToken, h.Profile)

	api.Get("/files", h.ReadToken, h.GetFiles)
	api.Post("/files", h.ReadToken, h.UploadFile)
	api.Get("/files/:id", h.ReadToken, h.DownloadFile)
	api.Delete("/files/:id", h.ReadToken, h.DeleteFile)

}
