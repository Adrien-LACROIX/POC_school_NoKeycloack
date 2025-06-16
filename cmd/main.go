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

func loadEnv() {
	paths := []string{
		"../configs/.env", // pour exécution depuis /cmd
		"./configs/.env",  // utile si tu lances go run depuis la racine
		".env",            // fallback racine projet
	}

	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Println("Fichier .env chargé :", path)
			return
		}
	}

	log.Fatal("Impossible de charger un fichier .env")
}
