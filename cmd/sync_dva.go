package main

import (
	"context"
	"log"

	"github.com/Chrisentech/aluta-market-api/database"
	"github.com/Chrisentech/aluta-market-api/internals/store"
)

func main() {
	// Initialize database connection
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	// Create store repository and service
	repo := store.NewRepository()
	svc := store.NewService(repo)

	// Create context
	ctx := context.Background()

	// Trigger sync
	log.Println("Starting Paystack DVA accounts sync...")
	if err := svc.SyncExistingPaystackDVAAccounts(ctx); err != nil {
		log.Fatalf("Failed to sync Paystack DVA accounts: %v", err)
	}
	log.Println("Successfully completed Paystack DVA accounts sync")
}
