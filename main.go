package main

import (
	"go-cloud/database"
	"go-cloud/handler"
	"go-cloud/repository"
	"go-cloud/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	app := fiber.New()

	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	serv := services.NewServices(repo)
	handler := handler.NewHandler(serv)

	startRoutingApp(app, handler)

	//3000 не использовать, для фронтэнда
	app.Listen(":8080")
}
