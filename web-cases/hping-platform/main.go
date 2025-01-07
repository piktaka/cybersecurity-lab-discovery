package main

import (
	"fmt"
	"net/http"

	"lablabee.com/cybersecurity-discovery1/hping-plateform/database"
	"lablabee.com/cybersecurity-discovery1/hping-plateform/handlers"
	"lablabee.com/cybersecurity-discovery1/hping-plateform/models"
)

func main() {
	database.Initialize()
	models.Migrate()
	fmt.Println("before handlers")
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/authenticate", handlers.HandleLogin)
	http.HandleFunc("/feed", handlers.FeedPage)
	http.HandleFunc("/feed/posts", handlers.HandlePostCreation)

	http.HandleFunc("/feed/{postId}/comments", handlers.HandleCommentCreation)
	fmt.Println("after handlers")

	// http.HandleFunc("/feed/likes", handlers.FeedPage)
	// http.HandleFunc("/feed/comments/:id", handlers.FeedPage)
	models.InsertUser("lablabee", "lablabee#2025@!")
	models.CreatePost("Hello This is the first post in the hping part of the lab")
	models.CreatePost("Hello This another post in the hping part of the lab Hehe")
	models.CreatePost("Can you stop the platform from working !! Hehe")

	http.ListenAndServe(":8080", nil)

}
