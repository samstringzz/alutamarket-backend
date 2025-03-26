package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize and start the server
	if err := InitServer(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
