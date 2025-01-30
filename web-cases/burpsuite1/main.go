package main

import (
	"net/http"

	"lablabee.com/cybersecurity-discovery1/case1/database"
	"lablabee.com/cybersecurity-discovery1/case1/handlers"
	"lablabee.com/cybersecurity-discovery1/case1/models"
)

func main() {
	database.Initialize()
	models.Migrate()
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/authenticate", handlers.HandleLogin)
	http.HandleFunc("/home", handlers.HomePage)
	http.HandleFunc("/logout", handlers.Logout)

	models.InsertUser("lablabee", "lablabee123")
	http.ListenAndServe(":8080", nil)

}
