package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"lablabee.com/cybersecurity-discovery1/hping-plateform/models"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

type User struct {
	Username string
	Password string
}

type Post struct {
	ID      uint      `gorm:"primaryKey;autoIncrement"`
	Content string    `gorm:"type:text;not null"`
	Comment []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
type Comment struct {
	ID      uint   `json:"id"`
	PostID  uint   `json:"post_id"`
	Comment string `json:"comment"`
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

	username := session.Values["username"].(string)

	temp := template.Must(template.ParseFiles("feed.html"))
	temp.Execute(w, struct{ Username string }{Username: username})

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

		vars := mux.Vars(r)
		postIdString := vars["postId"]
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
		databaseComment, err := models.CreateComment(uint(postId), comment.Comment)
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
