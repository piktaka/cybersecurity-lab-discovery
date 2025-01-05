package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {
	var err error

	DB, err = gorm.Open(sqlite.Open("user.db"), &gorm.Config{})

	if err != nil {

		log.Fatalf("Enable to create or get database: %v", err)

	}
	log.Println("Database connected successfully")

}
