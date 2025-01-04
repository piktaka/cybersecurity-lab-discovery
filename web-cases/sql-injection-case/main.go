package main

import (
	"net/http"

	"lablabee.com/cybersecurity-discovery1/sql-injection/database"
	"lablabee.com/cybersecurity-discovery1/sql-injection/handlers"
	"lablabee.com/cybersecurity-discovery1/sql-injection/models"
)

func main() {
	database.InitializeDB()
	database.InitializeNormalDB()
	models.Migrate()
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/authenticate", handlers.HandleLogin)
	http.HandleFunc("/home", handlers.HomePage)
	models.InsertUser("pikta", "pikta")
	http.ListenAndServe(":8080", nil)

}
