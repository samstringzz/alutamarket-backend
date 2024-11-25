package db

import (
	"log"
	"os"

	"github.com/Chrisentech/aluta-market-api/internals/cart"
	"github.com/Chrisentech/aluta-market-api/internals/messages"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/skynet"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Migrate() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbURI := os.Getenv("DB_URI")

	// Initialize the database with logging set to Silent
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	// Auto-migrate models to create tables if they don't exist
	db.AutoMigrate(&store.Store{}, &store.Invoice{}, &user.User{}, &product.Product{},
		&cart.Cart{}, &product.Category{}, &messages.Chat{}, &product.HandledProduct{},
		&store.Review{},
		&messages.Message{}, &store.Order{}, &store.Downloads{}, &user.PasswordReset{}, &skynet.Skynet{})

	return db
}
