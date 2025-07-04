// handler package allowed the navigation on the website
package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	db       *sql.DB
	tmpl     *template.Template
	tryCount = make(map[string]int)
)

// Handler connects to the database, defines the server handlers and listens for events on localhost:8080
func Handler() {
	var err error

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	tmpl = template.Must(template.ParseGlob("../web/templates/*.html"))

	http.Handle("/static/", http.StripPrefix("../web/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/signup", SignupHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/welcome", WelcomeHandler)
	http.HandleFunc("/trylater", TryLaterHandler)
	http.HandleFunc("/logout", LogoutHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Default handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

// Checks that no account exists with the same email address or username.
// Also checks that the passwords entered are identical.
// Saves the new user in the database by hashing the password.
//
// Redirects to the login page.
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "signup.html", nil)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	pwd := r.FormValue("pwd")
	pwd2 := r.FormValue("pwd2")

	data := map[string]string{}

	if email == "" || username == "" || pwd == "" || pwd2 == "" {
		data["Error"] = "Tous les champs sont requis"
		tmpl.ExecuteTemplate(w, "signup.html", data)
		return
	}
	if pwd != pwd2 {
		data["Error"] = "Les mots de passe ne correspondent pas"
		tmpl.ExecuteTemplate(w, "signup.html", data)
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 OR username=$2)", email, username).Scan(&exists)
	if err != nil {
		http.Error(w, "Erreur serveur", 500)
		return
	}
	if exists {
		data["Error"] = "Email ou pseudo déjà utilisé"
		tmpl.ExecuteTemplate(w, "signup.html", data)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur serveur", 500)
		return
	}
	_, err = db.Exec("INSERT INTO users(email, username, password) VALUES($1,$2,$3)", email, username, string(hash))
	if err != nil {
		http.Error(w, "Erreur lors de l'enregistrement", 500)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Verify that the user exists and compares password hashes.
// After three failed attempts, redirects to the "try later" page.
// If the login is successful, redirects to the "welcome" page.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	username := r.FormValue("username")
	pwd := r.FormValue("pwd")
	ip := r.RemoteAddr

	if tryCount[ip] >= 3 {
		http.Redirect(w, r, "/trylater", http.StatusSeeOther)
		return
	}

	var storedHash, dbUsername string
	err := db.QueryRow("SELECT password, username FROM users WHERE username=$1", username).Scan(&storedHash, &dbUsername)
	if err != nil {
		tryCount[ip]++
		tmpl.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Utilisateur/mot de passe incorrect"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(pwd))
	if err != nil {
		tryCount[ip]++
		tmpl.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Utilisateur/mot de passe incorrect"})
		return
	}

	tryCount[ip] = 0

	cookie := &http.Cookie{
		Name:    "session",
		Value:   dbUsername,
		Expires: time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

// displays "Welcome" followed by the user's name
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "welcome.html", map[string]string{"Username": cookie.Value})
}

// displays "Try later"
func TryLaterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "try_later.html", nil)
}

// logs the user out and returns them to the home page
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Secure:   true, // <-- Ajouté pour SonarQube
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
