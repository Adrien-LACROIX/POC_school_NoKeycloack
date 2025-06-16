package main

import (
	"log"

	"github.com/joho/godotenv"

	h "pocnokc/internal/handler"
)

func main() {
	loadEnv()

	h.Handler()
}

// loadEnv traverses a set of paths to find the .env file
// It loads the .env file once it is found.
// Fatal error on failure.
func loadEnv() {
	paths := []string{
		"../configs/.env", // for execution from /cmd
		"./configs/.env",  // useful if go run launched from root
		".env",            // fallback root project
	}

	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Println("Fichier .env chargé :", path)
			return
		}
	}

	log.Fatal("Impossible de charger un fichier .env")
}
