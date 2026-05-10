package main

import (
	"github.com/gofiber/fiber/v3"
)

func main() {

	app := fiber.New()

	startRoutingApp(app)

	app.Listen(":3000")
}
