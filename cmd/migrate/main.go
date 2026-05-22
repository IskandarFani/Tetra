package main

import (
	"go-cloud/database"
	"go-cloud/models"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("migrations completed successfully")
}
