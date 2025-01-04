package models

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm/clause"
	"lablabee.com/cybersecurity-discovery1/sql-injection/database"
)

type User struct {
	ID       uint   
	Username string 
	Password string 
}

func Migrate() {

	if err := database.DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to migrate User model: %v", err)
	}

	log.Println("User model migrated successfully")
}

func InsertUser(username, password string) error {
	if err := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&User{Username: username, Password: password}).Error; err != nil {
		return err
	}
	return nil
}

func GetUser(username,password string)(error) {
query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s' AND password='%s'", username,password)
			fmt.Println("Executing query:", query)
			row := database.NormalDB.QueryRow(query)
			user:=&User{}
			err:=row.Scan(&user.ID,&user.Username,&user.Password)
return err
}
