package main

import (
	"go-cloud/database"
	"go-cloud/repository"
	"go-cloud/services"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {

	app := fiber.New()

	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	serv := services.NewServices(repo)

	startRoutingApp(app)

	//3000 не использовать, для фронтэнда
	app.Listen(":8080")
}
