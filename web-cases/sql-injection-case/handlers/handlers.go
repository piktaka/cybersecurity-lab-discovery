package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"lablabee.com/cybersecurity-discovery1/sql-injection/models"
)

var store = sessions.NewCookieStore([]byte("super-secret-key-2"))

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
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
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
	session, _ := store.Get(r, "session")
	authenticated, ok := session.Values["authenticated"].(bool)
	if ok && authenticated {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()

	if err != nil {
		panic(err)
	}

	user := User{}

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	userfromDB, err := models.GetUser(user.Username, user.Password)

	if err != nil {

		if err == sql.ErrNoRows {
			renderLoginPageWithError(w, "Error: user or password incorrect")
			return
		} else {
			renderLoginPageWithError(w, "Error: Something went wrong")
			return
		}

	}
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	session, _ = store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = userfromDB.Username
	session.Save(r, w)
if user.Username == "lablabee'--" {
		filePath := filepath.Join("/opt/validation", "validate1")
		content := "challenge1"
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// File does not exist, create it and write the content
			err = os.WriteFile(filePath, []byte(content), 0644) // Read/write for owner, read for others
			if err != nil {
				fmt.Println("Error creating/writing file:", err)
				// Log the error appropriately.
			} else {
				fmt.Println("File /opt/validation/validate1 created with content:", content)
				// Optionally, log the successful creation.
			}
		} else if err != nil {
			// Error checking if the file exists
			fmt.Println("Error checking if file exists:", err)
			// Log this error.
		} else {
			fmt.Println("File /opt/validation/validate1 already exists.")
			// Optionally, handle the case where the file already exists (e.g., log it).
		}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
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

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	// Clear session data
	session.Values["authenticated"] = false
	session.Values["username"] = nil
	session.Save(r, w)
	// Redirect to the login page after logout
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
