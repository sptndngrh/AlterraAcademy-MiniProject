package main

import (
	"log"
	"sewakeun_project/configs"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	routes := configs.SetupRoutes()
	err := routes.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
