package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	// "lablabee.com/cybersecurity-discovery1/hping-plateform/models"
	"hping-platform/models"
)

var store = sessions.NewCookieStore([]byte("super-secret-key-3"))

type User struct {
	Username string
	Password string
}

type Post struct {
	ID       uint      `json:"id"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
}
type Comment struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	Comment   string    `json:"comment"`
	Timestamp time.Time `json:"timestamp"`
}

// func HomePage(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "session")
// 	authenticated, ok := session.Values["authenticated"].(bool)
// 	if !ok || !authenticated {
// 		http.Error(w, "You are not authenticated", http.StatusInternalServerError)
// 		return
// 	}

// 	username := session.Values["username"].(string)

// 	temp := template.Must(template.ParseFiles("home.html"))
// 	temp.Execute(w, struct{ Username string }{Username: username})

// }
func FeedPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		http.Error(w, "You are not authenticated", http.StatusInternalServerError)
		return
	}

	// username := session.Values["username"].(string)

	temp := template.Must(template.ParseFiles("feed.html"))
	posts, err := models.GetAllPosts()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	temp.Execute(w, posts)

}
func LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	session, _ := store.Get(r, "session")
	authenticated, ok := session.Values["authenticated"].(bool)
	if ok && authenticated {
		http.Redirect(w, r, "/feed", http.StatusSeeOther)
		return
	}
	temp := template.Must(template.ParseFiles("login.html"))
	temp.Execute(w, map[string]interface{}{
		"Error": "",
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()

	if err != nil {
		panic(err)
	}

	user := User{}

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	userFromDatabase, err := models.GetUser(user.Username)

	if err != nil {
		renderLoginPageWithError(w, "Error: user or password incorrect")
		return
	}

	if userFromDatabase.Password != user.Password {
		renderLoginPageWithError(w, "Error: user or password incorrect")
		return
	}
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = userFromDatabase.Username
	session.Save(r, w)
	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}

func renderLoginPageWithError(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusBadRequest)
	temp := template.Must(template.ParseFiles("login.html"))
	temp.Execute(w, map[string]interface{}{
		"Error": errorMessage,
	})
}

func HandleCommentCreation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		path := r.URL.Path
		parts := strings.Split(path, "/")

		// postIdString := vars["postId"]
		postIdString := parts[2]
		fmt.Println(postIdString)
		postId, err := strconv.ParseUint(postIdString, 10, 64)
		fmt.Println("this is the postID", postId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		comment := Comment{}
		bodyDecoder := json.NewDecoder(r.Body)
		err = bodyDecoder.Decode(&comment)
		comment.PostID = uint(postId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		comment.Timestamp = time.Now()
		databaseComment, err := models.CreateComment(uint(postId), comment.Comment, comment.Timestamp)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		comment.ID = databaseComment.ID
		_, err = models.AddComment(*databaseComment)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		bodyEncoder := json.NewEncoder(w)
		err = bodyEncoder.Encode(comment)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		w.WriteHeader(http.StatusCreated)

	}
}

func HandlePostCreation(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		post := Post{}
		bodyDecoder := json.NewDecoder(r.Body)
		err := bodyDecoder.Decode(&post)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		postFromDatabase, err := models.CreatePost(post.Content)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		post.ID = postFromDatabase.ID
		bodyEncoder := json.NewEncoder(w)
		err = bodyEncoder.Encode(post)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	if r.Method == http.MethodGet {
		postsFromDB, err := models.GetAllPosts()
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		posts := PostsfromDBtoJSON(postsFromDB)
		log.Println(posts)
		bodyEncoder := json.NewEncoder(w)
		err = bodyEncoder.Encode(posts)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		w.WriteHeader(http.StatusOK)

	}

}

func PostsfromDBtoJSON(dbPosts []models.Post) []Post {
	var post Post
	posts := make([]Post, 0, len(dbPosts))
	for _, dbPost := range dbPosts {
		post = Post{ID: dbPost.ID, Content: dbPost.Content, Comments: CommentsfromDBtoJSON(dbPost.Comments)}
		posts = append(posts, post)
	}
	return posts
}

func CommentsfromDBtoJSON(dbComments []models.Comment) []Comment {
	var comment Comment
	comments := make([]Comment, 0, len(dbComments))
	for _, dbComment := range dbComments {
		comment = Comment{ID: dbComment.ID, Comment: dbComment.Content, Timestamp: dbComment.Timestamp}
		comments = append(comments, comment)
	}
	return comments

}
