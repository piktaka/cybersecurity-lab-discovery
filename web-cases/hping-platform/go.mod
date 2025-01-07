// module lablabee.com/cybersecurity-discovery1/hping-plateform

module hping-platform

go 1.23.4

// replace lablabee.com/cybersecurity-discovery1/hping-plateform/database => ./database

// replace lablabee.com/cybersecurity-discovery1/hping-plateform/handlers => ./handlers

// replace lablabee.com/cybersecurity-discovery1/hping-plateform/models => ./models

require (
	github.com/gorilla/sessions v1.4.0
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	golang.org/x/text v0.21.0 // indirect
)
