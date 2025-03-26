package utils

import (
	"gorm.io/gorm"
	"log"
	"time"
)

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
