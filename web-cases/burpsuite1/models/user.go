package models

import (
	"log"

	"gorm.io/gorm/clause"
	"lablabee.com/cybersecurity-discovery1/case1/database"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
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

func GetUser(username string) (*User, error) {
	user := &User{}
	result := database.DB.Where("username = ?", username).First(user)
	if err := result.Error; err != nil {
		return nil, err

	}
	return user, nil

}
