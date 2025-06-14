package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	db       *sql.DB
	tmpl     *template.Template
	tryCount = make(map[string]int)
)

func main() {
	var err error
	connStr := "host=localhost port=5432 user=myuser password=mypass dbname=myappdb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email TEXT UNIQUE NOT NULL,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`)
	if err != nil {
		log.Fatal(err)
	}

	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/try-later", tryLaterHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	username := r.FormValue("username")
	pwd := r.FormValue("pwd")
	ip := r.RemoteAddr

	if tryCount[ip] >= 3 {
		http.Redirect(w, r, "/try-later", http.StatusSeeOther)
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

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "welcome.html", map[string]string{"Username": cookie.Value})
}

func tryLaterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "try_later.html", nil)
}
