package database

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		dbURI := os.Getenv("DATABASE_URL")
		if dbURI == "" {
			log.Fatal("DATABASE_URL environment variable not set")
		}

		db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
			PrepareStmt: true,
		})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get database instance: %v", err)
		}

		sqlDB.SetMaxIdleConns(2)
		sqlDB.SetMaxOpenConns(5)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)

		instance = db
	})
	return instance
}
