package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Fatal("DB_URI environment variable is not set")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Error creating migration driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Error creating migration instance:", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Please specify 'up' or 'down'")
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Error running migrations:", err)
		}
		fmt.Println("Successfully ran migrations up")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Error running migrations:", err)
		}
		fmt.Println("Successfully ran migrations down")
	default:
		log.Fatal("Invalid command. Please use 'up' or 'down'")
	}
}
