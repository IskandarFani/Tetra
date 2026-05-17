package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {

	dsn := "host=127.0.0.1 port=5432 user=database_user password=database_password dbname=database_name sslmode=disable" //err

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil

}
