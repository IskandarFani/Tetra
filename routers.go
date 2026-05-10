package main

import (
	"go-cloud/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func startRoutingApp(app *fiber.App) {
	app.Get("/", static.New("./public/hello.html"))
	app.Get("/health", handler.Health())
	app.Get("/monitor", handler.Monitor())
	app.Get("/checkdb", handler.CheckDB())
}
