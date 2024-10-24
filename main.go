package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/noobpiyush/paytm-api/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	log.Print(PORT)

	routes.RegisteredRoutes()
	log.Printf("Server connected on port 8080...")
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
