package db

import (
	"database/sql"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabse() (*Database, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URI")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
