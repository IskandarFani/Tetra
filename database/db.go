package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DBStatus struct {
	Connected bool
	Message   string
}

func Connect() DBStatus {

	dsn := "host=127.0.0.1 port=5432 user=database_user password=database_password dbname=database_name sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return DBStatus{
			Connected: false,
			Message:   fmt.Sprintf("Failed to connect database: %v", err),
		}
	}

	DB = db

	return DBStatus{
		Connected: true,
		Message:   "Connected to Postgres successfully",
	}

}
