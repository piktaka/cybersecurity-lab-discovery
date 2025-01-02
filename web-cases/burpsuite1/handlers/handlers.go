package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"lablabee.com/cybersecurity-discovery1/case1/models"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

type User struct {
	Username string
	Password string
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		http.Error(w, "You are not authenticated", http.StatusInternalServerError)
		return
	}

	username := session.Values["username"].(string)

	temp := template.Must(template.ParseFiles("home.html"))
	temp.Execute(w, struct{ Username string }{Username: username})

}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	authenticated, ok := session.Values["authenticated"].(bool)
	if ok && authenticated {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
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

	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = userFromDatabase.Username
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusCreated)
}

func renderLoginPageWithError(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusBadRequest)
	temp := template.Must(template.ParseFiles("login.html"))
	temp.Execute(w, map[string]interface{}{
		"Error": errorMessage,
	})
}
