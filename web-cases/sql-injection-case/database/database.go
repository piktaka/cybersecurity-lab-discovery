package database

import (
	"database/sql"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var NormalDB *sql.DB
func InitializeDB() {
	var err error

	DB, err = gorm.Open(sqlite.Open("user.db"), &gorm.Config{})

	if err != nil {

		log.Fatalf("Enable to create or get database: %v", err)

	}
	log.Println("Database connected successfully")


}

func InitializeNormalDB(){
	var err error

		NormalDB, err = sql.Open("sqlite3", "user.db")
	if err!=nil {
		log.Fatal(err)
	}
		
}
