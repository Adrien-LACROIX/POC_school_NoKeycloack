package main

import (
	"database/sql"
	"log"

	h "pocnokc/internal/handler"
)

var db *sql.DB

func main() {
	err := launchDB()
	if err != nil {
		log.Fatal(err)
	}

	connStr := "host=localhost port=5432 user=postgres password=mypwd dbname=myappdb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	h.Handler(db)
}
