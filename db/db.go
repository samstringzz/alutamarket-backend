package db

import (
	"github.com/Chrisentech/aluta-market-api/database"
	"github.com/Chrisentech/aluta-market-api/internals/cart"
	"github.com/Chrisentech/aluta-market-api/internals/messages"
	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/skynet"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate() *gorm.DB {
	db := database.GetDB() // Use the database package's GetDB function

	// Auto-migrate models to create tables if they don't exist
	if err := db.AutoMigrate(
		&store.Store{},
		&store.Invoice{},
		&user.User{},
		&product.Product{},
		&cart.Cart{},
		&product.Category{},
		&product.SubCategory{}, // Added SubCategory
		&messages.Chat{},
		&product.HandledProduct{},
		&store.Review{},
		&messages.Message{},
		&store.Order{},
		&store.Downloads{},
		&user.PasswordReset{},
		&skynet.Skynet{},
		&store.DVAAccount{},
		&store.DVACustomer{},
		&store.DVABank{},
	); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	return db
}
