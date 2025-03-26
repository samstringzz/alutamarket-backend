package utils

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// RetryConnection attempts to establish a database connection with retries
func RetryConnection(db *gorm.DB, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Failed to get database instance (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		err = sqlDB.Ping()
		if err == nil {
			return nil
		}

		log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return err
}

// ConfigureConnectionPool sets up the database connection pool with recommended settings
func ConfigureConnectionPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return nil
}

// CheckConnection verifies if the database connection is alive
func CheckConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
