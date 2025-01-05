package main

import (
	"net/http"

	"lablabee.com/cybersecurity-discovery1/hping-plateform/database"
	"lablabee.com/cybersecurity-discovery1/hping-plateform/handlers"
	"lablabee.com/cybersecurity-discovery1/hping-plateform/models"
)

func main() {
	database.Initialize()
	models.Migrate()
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/authenticate", handlers.HandleLogin)
	http.HandleFunc("/home", handlers.HomePage)
	models.InsertUser("pikta", "pikta")
	http.ListenAndServe(":8080", nil)

}
