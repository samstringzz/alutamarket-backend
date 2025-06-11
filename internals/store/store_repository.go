package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samstringzz/alutamarket-backend/database"
	"github.com/samstringzz/alutamarket-backend/errors"
	"github.com/samstringzz/alutamarket-backend/internals/paystack"
	"github.com/samstringzz/alutamarket-backend/utils"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

func (r *repository) CheckStoreName(ctx context.Context, query string) error {
	var stores []*Store
	if err := r.db.Where("name ILIKE ?", "%"+query+"%").Find(&stores).Error; err != nil {
		return err
	}
	for _, item := range stores {
		if item.Name == query {
			return errors.NewAppError(http.StatusConflict, "CONFLICT", "Store Name already choosen")
		}
	}
	return nil
}

func (r *repository) CreateStore(ctx context.Context, req *Store) (*Store, error) {
	var count int64
	r.db.Model(&Store{}).Where("name = ? AND user_id = ?", req.Name, req.UserID).Count(&count)
	if count > 0 {
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Store already exists")
	}

	resp := &Store{
		Name:               req.Name,
		Email:              req.Email,
		Link:               req.Link,
		UserID:             req.UserID,
		Description:        req.Description,
		HasPhysicalAddress: req.HasPhysicalAddress,
		Address:            req.Address,
		Wallet:             0,
		Status:             true,
		Background:         "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQQLbvWGTFQh6OGWPfkLx2xBS_OP3oZJzQubA&s",
		Phone:              req.Phone,
	}
	if err := r.db.Create(resp).Error; err != nil {
		r.db.Rollback()
		return nil, err
	}

	// Create Paystack DVA account
	paystackAccount, err := r.getPaystackDVAAccount(req.Email)
	if err != nil {
		log.Printf("Warning: Failed to create Paystack DVA account: %v", err)
		return resp, nil // Return store even if DVA creation fails
	}

	// Store the Paystack DVA account details
	if err := r.CreatePaystackDVAAccount(ctx, resp.ID, paystackAccount, req.Email); err != nil {
		log.Printf("Warning: Failed to store Paystack DVA account: %v", err)
		return resp, nil // Return store even if DVA storage fails
	}

	return resp, nil
}

func (r *repository) CreateInvoice(ctx context.Context, req *Invoice) (*Invoice, error) {
	_, err := r.GetStore(ctx, req.StoreID)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

func (r *repository) GetStore(ctx context.Context, id uint32) (*Store, error) {
	var store *Store
	err := r.db.Where("id = ?", id).First(&store).Error
	if err != nil {
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Store does not exist")
	}
	return store, nil
}

func (r *repository) GetStoreByName(ctx context.Context, name string) (*Store, error) {
	var store *Store
	err := r.db.Where("name = ? or link = ?", name, name).First(&store).Error
	if err != nil {
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Store does not exist")
	}
	return store, nil
}

func (r *repository) GetStores(ctx context.Context, userID uint32, limit int, offset int) ([]*Store, error) {
	var stores []*Store

	// Create the base query
	query := r.db.Table("stores")

	// Only add userID filter if it's not 0 (which means fetch all stores)
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	// Add pagination if limit is greater than 0
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	// Execute the query
	result := query.Find(&stores)
	if result.Error != nil {
		return nil, result.Error
	}

	// For each store, fetch related data
	for _, store := range stores {
		// Fetch followers with error handling
		if err := r.db.Table("store_followers").Where("store_id = ?", store.ID).Find(&store.Followers).Error; err != nil {
			log.Printf("Warning: Error fetching followers for store %d: %v", store.ID, err)
			store.Followers = []*Follower{} // Initialize empty array instead of nil
		}

		// Fetch products with error handling
		var products []Product
		if err := r.db.Table("products").
			Where("store = ? AND deleted_at IS NULL", store.Name). // Add store name condition and check for soft deletes
			Find(&products).Error; err != nil {
			log.Printf("Warning: Error fetching products for store %d: %v", store.ID, err)
			store.Products = []Product{} // Initialize empty array instead of nil
		} else {
			store.Products = products
		}

		// Fetch accounts with error handling
		if err := r.db.Table("dva_accounts").Where("store_id = ?", store.ID).Find(&store.Accounts).Error; err != nil {
			log.Printf("Warning: Error fetching accounts for store %d: %v", store.ID, err)
			store.Accounts = []*WithdrawalAccount{} // Initialize empty array instead of nil
		}

		// Initialize empty arrays if they're nil
		if store.Visitors == nil {
			store.Visitors = []string{}
		}

		// Fetch orders for this store with proper error handling
		var orders []*Order
		if err := r.db.Table("orders").
			Where("? = ANY(stores_id)", store.Name).
			Find(&orders).Error; err != nil {
			log.Printf("Warning: Error fetching orders for store %d: %v", store.ID, err)
			store.Orders = []*StoreOrder{} // Initialize empty array instead of nil
		} else {
			// Convert orders to store orders with proper error handling
			var storeOrders []*StoreOrder
			for _, order := range orders {
				// Skip invalid orders
				if order == nil {
					continue
				}

				var details DeliveryDetails
				var customer Customer
				var storeProducts []*StoreProduct

				// Safely unmarshal delivery details
				if order.DeliveryDetailsJSON != "" {
					if err := json.Unmarshal([]byte(order.DeliveryDetailsJSON), &details); err != nil {
						log.Printf("Warning: Error unmarshaling delivery details for order %s: %v", order.UUID, err)
						continue
					}
					order.DeliveryDetails = &details
				}

				// Safely unmarshal customer details
				if order.CustomerJSON != "" {
					if err := json.Unmarshal([]byte(order.CustomerJSON), &customer); err != nil {
						log.Printf("Warning: Error unmarshaling customer for order %s: %v", order.UUID, err)
						continue
					}
					order.Customer = &customer
				}

				// Skip if required data is missing
				if order.Customer == nil {
					log.Printf("Warning: Missing customer data for order %s", order.UUID)
					continue
				}

				// Fetch product details for each order
				if len(order.Products) > 0 {
					for _, p := range order.Products {
						if p.ID == 0 {
							continue // Skip invalid products
						}

						// Get full product details from products table
						var fullProduct Product
						if err := r.db.Table("products").Where("id = ?", p.ID).First(&fullProduct).Error; err != nil {
							log.Printf("Warning: Error fetching product details for ID %d: %v", p.ID, err)
							continue
						}

						storeProducts = append(storeProducts, &StoreProduct{
							ID:        fullProduct.ID,
							Name:      fullProduct.Name,
							Thumbnail: fullProduct.Thumbnail,
							Price:     fullProduct.Price,
							Quantity:  p.Quantity, // Use quantity from order
							Status:    "active",   // Default to active since it was ordered
						})
					}
				}

				storeOrder := &StoreOrder{
					StoreID:   strconv.FormatUint(uint64(store.ID), 10),
					Products:  storeProducts,
					Status:    order.Status,
					TransRef:  order.TransRef,
					UUID:      order.UUID,
					Customer:  *order.Customer,
					CreatedAt: order.CreatedAt,
					UpdatedAt: order.UpdatedAt,
				}
				storeOrders = append(storeOrders, storeOrder)

				// Debug log for each order
				log.Printf("Processing order %s with %d products for store %s",
					order.UUID,
					len(storeProducts),
					store.Name)
			}

			store.Orders = storeOrders
			// Debug log for store orders
			log.Printf("Store %s: Added %d orders", store.Name, len(storeOrders))
		}

		// Log the data being returned for debugging
		log.Printf("Store %s: %d products, %d orders", store.Name, len(store.Products), len(store.Orders))
	}

	return stores, nil
}

func (r *repository) GetInvoices(ctx context.Context, storeId uint32) ([]*Invoice, error) {
	var invoice []*Invoice
	if err := r.db.Where("store_id=?", storeId).Find(&invoice).Error; err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *repository) DeleteStore(ctx context.Context, id uint32) error {
	existingStore, err := r.GetStore(ctx, id)
	if err != nil {
		return err
	}
	err = r.db.Where("id = ?", id).First(&Store{}).Error
	if err != nil {
		return err
	}
	err = r.db.Delete(existingStore).Error
	return err
}

func (r *repository) UpdateStore(ctx context.Context, req *UpdateStore) (*Store, error) {
	// First, check if the Store exists by its ID or another unique identifier
	existingStore, err := r.GetStore(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Update only the fields that are present in the req
	if req.Name != "" {
		existingStore.Name = req.Name
	}
	if req.Description != "" {
		existingStore.Description = req.Description
	}
	if len(req.Visitors) > 0 {
		// Convert visitor IDs to strings and ensure no duplicates
		visitorMap := make(map[string]bool)
		for _, v := range existingStore.Visitors {
			visitorMap[v] = true
		}
		for _, v := range req.Visitors {
			visitorMap[v] = true
		}

		// Convert map back to slice
		existingStore.Visitors = make(pq.StringArray, 0, len(visitorMap))
		for v := range visitorMap {
			existingStore.Visitors = append(existingStore.Visitors, v)
		}
		sort.Strings(existingStore.Visitors)
	}
	if req.Link != "" {
		existingStore.Link = req.Link
	}
	if req.Phone != "" {
		existingStore.Phone = req.Phone
	}
	if req.Background != "" {
		existingStore.Background = req.Background
	}
	if req.Thumbnail != "" {
		existingStore.Thumbnail = req.Thumbnail
	}
	if req.Email != "" {
		existingStore.Email = req.Email
	}
	existingStore.Wallet += req.Wallet

	// Update maintenance mode if provided
	existingStore.MaintenanceMode = req.MaintenanceMode

	// Handle account update only if account information is provided
	if req.Account != nil {
		// Create a new account record directly in dva_accounts table
		bankDetails := map[string]interface{}{
			"store_id":       req.ID,
			"account_number": req.Account.AccountNumber,
			"account_name":   req.Account.AccountName,
			"bank_name":      req.Account.BankName,
			"bank_code":      req.Account.BankCode,
			"bank_image":     req.Account.BankImage,
			"bank_id":        req.Account.BankID,
		}

		// Try to update existing record first
		result := r.db.Table("dva_accounts").
			Where("store_id = ?", req.ID).
			Updates(bankDetails)

		if result.Error != nil {
			return nil, fmt.Errorf("failed to update bank details: %v", result.Error)
		}

		// If no record was updated, create a new one
		if result.RowsAffected == 0 {
			// Get the existing store to get user details
			existingStore, err := r.GetStore(ctx, req.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get store details: %v", err)
			}

			// Get the user details
			var user struct {
				ID       string `json:"id"`
				UUID     string `json:"uuid"`
				Fullname string `json:"fullname"`
				Email    string `json:"email"`
			}
			if err := r.db.Table("users").
				Where("id = ?", existingStore.UserID).
				First(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to get user details: %v", err)
			}

			// Add customer_id to bank details
			bankDetails["customer_id"] = user.UUID

			// Create new record
			if err := r.db.Table("dva_accounts").Create(bankDetails).Error; err != nil {
				return nil, fmt.Errorf("failed to create bank details: %v", err)
			}
		}

		// Fetch accounts using the correct column name
		var accounts []*WithdrawalAccount
		if err := r.db.Table("dva_accounts").Where("store_id = ?", req.ID).Find(&accounts).Error; err != nil {
			log.Printf("Error fetching accounts for store %d: %v", req.ID, err)
		}
		existingStore.Accounts = accounts
	}

	// Update the Store in the repository
	if err := r.db.Save(existingStore).Error; err != nil {
		return nil, err
	}

	// Fetch the updated store with accounts
	var updatedStore Store
	if err := r.db.First(&updatedStore, existingStore.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload store: %v", err)
	}

	// Manually fetch accounts
	if err := r.db.Table("dva_accounts").Where("store_id = ?", updatedStore.ID).Find(&updatedStore.Accounts).Error; err != nil {
		log.Printf("Error fetching accounts for store %d: %v", updatedStore.ID, err)
	}

	return &updatedStore, nil
}

func (r *repository) CreateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error) {
	var store Store
	err := r.db.First(&store, req.StoreID).Error
	if err != nil {
		return nil, err
	}

	panic("not implementd")
}

func (r *repository) GetOrders(ctx context.Context, storeID uint32) ([]*Order, error) {
	var orders []*Order
	err := r.db.WithContext(ctx).Where("store_id = ?", storeID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *repository) GetOrdersByStore(ctx context.Context, storeName string) ([]*Order, error) {
	var orders []*Order

	// Query orders where the store name exists in the stores_id array
	err := r.db.
		Where("? = ANY(stores_id)", storeName).
		Find(&orders).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}

	// Unmarshal JSON fields for each order
	for _, order := range orders {
		// Unmarshal delivery details
		if order.DeliveryDetailsJSON != "" {
			var details DeliveryDetails
			if err := json.Unmarshal([]byte(order.DeliveryDetailsJSON), &details); err != nil {
				return nil, fmt.Errorf("failed to unmarshal delivery details: %v", err)
			}
			order.DeliveryDetails = &details
		}

		// Unmarshal customer details
		if order.CustomerJSON != "" {
			var customer Customer
			if err := json.Unmarshal([]byte(order.CustomerJSON), &customer); err != nil {
				return nil, fmt.Errorf("failed to unmarshal customer: %v", err)
			}
			order.Customer = &customer
		}
	}

	return orders, nil
}

func (r *repository) GetPurchasedOrders(ctx context.Context, userID string) ([]*Order, error) {
	var orders []*Order

	// Query orders with proper type conversion and preload relations
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if err := query.Find(&orders).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*Order{}, nil
		}
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}

	// Unmarshal JSON fields for each order
	for _, order := range orders {
		// Unmarshal delivery details
		if order.DeliveryDetailsJSON != "" {
			var details DeliveryDetails
			if err := json.Unmarshal([]byte(order.DeliveryDetailsJSON), &details); err != nil {
				return nil, fmt.Errorf("failed to unmarshal delivery details: %v", err)
			}
			order.DeliveryDetails = &details
		}

		// Unmarshal customer details
		if order.CustomerJSON != "" {
			var customer Customer
			if err := json.Unmarshal([]byte(order.CustomerJSON), &customer); err != nil {
				return nil, fmt.Errorf("failed to unmarshal customer: %v", err)
			}
			order.Customer = &customer
		}
	}

	return orders, nil
}

func (r *repository) UpdateOrderStatus(ctx context.Context, uuid string, status, transStatus string) error {
	// Get the order first
	order, err := r.GetOrderByUUID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to get order: %v", err)
	}

	// If changing from delivered to another status, remove the earnings
	if order.Status == "delivered" && status != "delivered" {
		// Find and update store earnings to mark as reversed
		var earnings StoreEarnings
		if err := r.db.Where("order_id = ? AND status = ?", uuid, "released").
			First(&earnings).Error; err == nil {
			earnings.Status = "reversed"
			if err := r.db.Save(&earnings).Error; err != nil {
				return fmt.Errorf("failed to reverse store earnings: %v", err)
			}
		}
	}

	// Update the order status
	if err := r.db.Model(&Order{}).
		Where("uuid = ?", uuid).
		Updates(map[string]interface{}{
			"status":       status,
			"trans_status": transStatus,
		}).Error; err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	// If the order is being marked as delivered, add to store earnings
	if status == "delivered" {
		amount, err := strconv.ParseFloat(order.Amount, 64)
		if err != nil {
			return fmt.Errorf("failed to parse order amount: %v", err)
		}

		// Check if earnings already exist for this order with status 'released'
		var existingEarnings StoreEarnings
		notFoundErr := r.db.Where("order_id = ? AND status = ?", uuid, "released").First(&existingEarnings).Error

		if notFoundErr != nil && notFoundErr != gorm.ErrRecordNotFound {
			// Handle other potential errors during the check
			log.Printf("Error checking for existing store earnings for order %s: %v", uuid, notFoundErr)
			return fmt.Errorf("failed to check for existing store earnings: %v", notFoundErr)
		} else if notFoundErr == gorm.ErrRecordNotFound {
			// Log that no existing earnings were found, so new ones will be created
			log.Printf("No existing released store earnings found for order %s. Proceeding to create.", uuid)

			// Log the store names and amount before attempting to create earnings
			log.Printf("Creating earnings for stores (by name): %v, order amount: %f", order.StoresID, amount)

			// Add earnings for each store in the order
			for _, storeName := range order.StoresID {
				// Fetch the store by name to get its ID
				var store Store
				if err := r.db.Where("name = ?", storeName).First(&store).Error; err != nil {
					log.Printf("Error finding store with name %s for order %s: %v. Skipping earnings creation for this store.", storeName, uuid, err)
					continue // Skip this store if not found
				}

				// Log the retrieved store ID
				log.Printf("Retrieved Store ID %d for store name %s", store.ID, storeName)

				earnings := &StoreEarnings{
					StoreID:   store.ID, // Use the retrieved numeric store ID
					OrderID:   uuid,
					Amount:    amount,
					Status:    "released",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				// Log before attempting to add store earnings
				log.Printf("AddStoreEarnings: Attempting to create earnings with StoreID=%d, OrderID=%s, Amount=%f, Status=%s", earnings.StoreID, earnings.OrderID, earnings.Amount, earnings.Status)

				if err := r.AddStoreEarnings(ctx, earnings); err != nil {
					// Log the error if adding store earnings fails
					log.Printf("AddStoreEarnings: Failed to create earnings: %v", err)
					// Do NOT return here, allow processing of other stores in the order
				}
			}
			// Log after the loop finishes
			log.Printf("Finished attempting to create earnings for order %s.", uuid)

		} else { // notFoundErr == nil
			// Log that existing released earnings were found, so no new ones will be created
			log.Printf("Existing released store earnings found for order %s. Skipping creation.", uuid)
		}
	}

	return nil
}

func (r *repository) GetOrder(ctx context.Context, storeID uint32, orderID string) (*Order, error) {
	var store Store
	err := r.db.WithContext(ctx).Preload("Orders", "id = ?", orderID).First(&store, storeID).Error
	if err != nil {
		return nil, err
	}

	panic("Not implented")

	// if len(store.Orders) == 0 {
	// 	return nil, gorm.ErrRecordNotFound
	// }

	// return store.Orders[0], nil
}

func (r *repository) GetFollowedStores(ctx context.Context, followerID uint32) ([]*Store, error) {
	var stores []*Store

	query := `
        SELECT * FROM stores 
        WHERE followers::jsonb @> '[{"follower_id": %d}]'
    `

	if err := r.db.Raw(fmt.Sprintf(query, followerID)).Scan(&stores).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch followed stores: %v", err)
	}

	return stores, nil
}

func (r *repository) UpdateOrder(ctx context.Context, req *UpdateStoreOrderInput) (*Order, error) {
	var existingOrder *Order
	var existingStore *Store

	// Fetch the store and order concurrently using goroutines
	errChan := make(chan error, 2) // Error channel to handle errors from goroutines

	go func() {
		errChan <- r.db.WithContext(ctx).Where("id = ?", req.StoreID).First(&existingStore).Error
	}()

	go func() {
		errChan <- r.db.WithContext(ctx).Where("uuid = ?", req.UUID).First(&existingOrder).Error
	}()

	// Wait for both operations to finish
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return nil, err // Return the first error encountered
		}
	}

	// Update the store's order status
	for _, storeOrder := range existingStore.Orders {
		if storeOrder.UUID == req.UUID {
			storeOrder.Status = req.Status
			storeOrder.UpdatedAt = time.Now()
		}
	}

	// Update the order fields
	existingOrder.Status = "In Progress"

	for _, product := range filterProductsByStore(existingOrder.Products, existingStore.Name) {
		product.Status = req.Status
	}
	existingOrder.UpdatedAt = time.Now()

	// Send email notifications concurrently
	go func() {
		to := []string{existingOrder.Customer.Email}
		contents := map[string]string{
			"buyer_name": existingOrder.Customer.Name,
			"order_id":   existingOrder.UUID,
		}

		if req.Status == "canceled" {
			templateID := "bb57c0b0-cb2b-4cd7-9170-f2c536a3dfe2"
			utils.SendEmail(templateID, "Your Order was Declined", to, contents)
		} else if req.Status == "processing" {
			templateID := "04551de0-1db2-46bb-b48a-610b744ee3e9"
			utils.SendEmail(templateID, "Your Order has been Confirmed", to, contents)
		}
	}()

	// Save changes to the database concurrently
	saveErrChan := make(chan error, 2)

	go func() {
		saveErrChan <- r.db.WithContext(ctx).Save(&existingStore).Error
	}()

	go func() {
		saveErrChan <- r.db.WithContext(ctx).Save(&existingOrder).Error
	}()

	// Wait for both save operations to finish
	for i := 0; i < 2; i++ {
		if err := <-saveErrChan; err != nil {
			return nil, err // Return the first error encountered
		}
	}

	return existingOrder, nil
}

func (r *repository) UpdateStoreFollowership(ctx context.Context, storeID uint32, follower *Follower, action string) (*Store, error) {
	// Log the incoming request
	log.Printf("Attempting to %s store with ID %d for follower ID %d", action, storeID, follower.FollowerID)

	// Fetch the store
	existingStore, err := r.GetStore(ctx, storeID)
	if err != nil {
		log.Printf("Error fetching store with ID %d: %v", storeID, err)
		return nil, err
	}

	// Check if the user is already a follower
	userExists := false
	var followerIndex int
	for i, existingFollower := range existingStore.Followers {
		if existingFollower.FollowerID == follower.FollowerID {
			userExists = true
			followerIndex = i
			break
		}
	}

	// Handle the follow/unfollow action
	switch action {
	case "follow":
		if userExists {
			// log.Printf("Follower with ID %d already exists in store %d", follower.FollowerID, storeID)
			return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "user is already a follower")
		}
		// log.Printf("Adding follower with ID %d to store %d", follower.FollowerID, storeID)
		existingStore.Followers = append(existingStore.Followers, follower)
	case "unfollow":
		if !userExists {
			// log.Printf("Follower with ID %d not found in store %d", follower.FollowerID, storeID)
			return nil, errors.NewAppError(http.StatusNotFound, "NOT FOUND", "user is not a follower")
		}
		// log.Printf("Removing follower with ID %d from store %d", follower.FollowerID, storeID)
		existingStore.Followers = append(existingStore.Followers[:followerIndex], existingStore.Followers[followerIndex+1:]...)
	default:
		// log.Printf("Invalid action: %s", action)
		return nil, errors.NewAppError(http.StatusNotAcceptable, "INVALID", "invalid action")
	}

	// Save the updated store
	err = r.db.WithContext(ctx).Save(&existingStore).Error
	if err != nil {
		// log.Printf("Error saving updated store with ID %d: %v", storeID, err)
		return nil, err
	}

	// log.Printf("Successfully updated store %d for action %s", storeID, action)
	return existingStore, nil
}

func (r *repository) GetStoresByFollower(ctx context.Context, followerID uint32) ([]*Store, error) {
	var stores []*Store
	if err := r.db.Where("followers_id=?", followerID).Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *repository) CreateTransactions(ctx context.Context, req *Transactions) (*Transactions, error) {

	var store *Store
	err := r.db.First(&store, req.StoreID).Error
	if err != nil {
		return nil, err
	}
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	req.UUID = "AM-" + utils.GenerateRandomString(6)

	store.Transactions = append(store.Transactions, req)

	err = r.db.WithContext(ctx).Save(&store).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *repository) WithdrawFund(ctx context.Context, req *Fund) error {
	// Create service
	service := NewService(r)

	// Process withdrawal
	return service.WithdrawFund(ctx, req)
}

// AddReview adds a new review to the database
func (r *repository) AddReview(ctx context.Context, review *Review) error {
	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		return err
	}
	return nil
}

// storeReviews, err := repo.GetReviews(ctx, "storeID", uint32(1))
// orderReviews, err := repo.GetReviews(ctx, "orderID", "order1
// GetReviews fetches reviews by either storeID or orderID based on the filter type.
func (r *repository) GetReviews(ctx context.Context, filterType string, value interface{}) ([]*Review, error) {
	var reviews []*Review

	// Determine the filter type and apply the query
	switch filterType {
	case "storeID":
		storeID, ok := value.(uint32)
		if !ok {
			return nil, errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "invalid storeID type")
		}
		if err := r.db.WithContext(ctx).Where("store_id = ?", storeID).Find(&reviews).Error; err != nil {
			return nil, err
		}
	case "productID":
		productID, ok := value.(uint32)
		if !ok {
			return nil, errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "invalid storeID type")
		}
		if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
			return nil, err
		}
	case "orderID":
		orderID, ok := value.(string)
		if !ok {
			return nil, errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "invalid orderID type")
		}
		if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&reviews).Error; err != nil {
			return nil, err
		}

	default:
		return nil, errors.NewAppError(http.StatusForbidden, "FORBIDDEN", "invalid filter type")
	}

	return reviews, nil
}

func filterProductsByStore(products []TrackedProduct, storeName string) []*StoreProduct {
	var filteredProducts []*StoreProduct

	for i, product := range products {
		product.Quantity = products[i].Quantity
		storeProduct := convertTrackedToStoreProduct(product)
		if product.Store == storeName {
			filteredProducts = append(filteredProducts, storeProduct)
		}
	}

	return filteredProducts
}

func convertTrackedToStoreProduct(tp TrackedProduct) *StoreProduct {
	return &StoreProduct{
		Name:      tp.Name,
		Price:     tp.Price,
		Quantity:  tp.Quantity,
		Thumbnail: tp.Thumbnail,
		ID:        tp.ID,
	}
}

func (r *repository) GetDVAAccount(ctx context.Context, email string) (*DVAAccount, error) {
	var account DVAAccount

	// First try to get from local database
	err := r.db.Table("dva_accounts").
		Select("dva_accounts.id, dva_accounts.customer_id, dva_accounts.bank_id, dva_accounts.account_number, dva_accounts.account_name").
		Preload("Customer").
		Preload("Bank").
		Joins("JOIN dva_customers ON dva_accounts.customer_id = dva_customers.id").
		Joins("JOIN dva_banks ON dva_accounts.bank_id = dva_banks.id").
		Where("dva_customers.email = ?", email).
		First(&account).Error

	if err != nil {
		// If not found in database, check Paystack
		paystackAccount, paystackErr := r.getPaystackDVAAccount(email)
		if paystackErr != nil {
			return nil, fmt.Errorf("account not found in database or Paystack: %v", paystackErr)
		}

		// Create string IDs with prefixes and random strings
		timestamp := time.Now().Unix()
		account = DVAAccount{
			ID:            fmt.Sprintf("DVA_%d_%s", timestamp, utils.GenerateRandomString(8)),
			AccountNumber: paystackAccount.AccountNumber,
			AccountName:   paystackAccount.AccountName,
			Customer: DVACustomer{
				ID:    fmt.Sprintf("CUST_%d_%s", timestamp, utils.GenerateRandomString(8)),
				Email: email,
			},
			Bank: DVABank{
				ID:   fmt.Sprintf("BANK_%d_%s", timestamp, utils.GenerateRandomString(8)),
				Name: paystackAccount.Bank.Name,
				Slug: "wema-bank",
			},
		}

		// Save to database
		if err := r.db.Create(&account).Error; err != nil {
			return nil, fmt.Errorf("failed to save Paystack account to database: %v", err)
		}
	}

	return &account, nil
}

func (r *repository) getPaystackDVAAccount(email string) (*PaystackDVAResponse, error) {
	// Create Paystack client
	paystackClient := paystack.NewClient(os.Getenv("PAYSTACK_SECRET_KEY"))

	// Get DVA account using the client
	account, err := paystackClient.GetDVAAccount(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get Paystack DVA account: %v", err)
	}

	// If no account found, return error
	if account == nil {
		return nil, fmt.Errorf("no DVA account found for email: %s", email)
	}

	// Convert to our response format
	return &PaystackDVAResponse{
		AccountNumber: account.AccountNumber,
		AccountName:   account.AccountName,
		Bank: struct {
			Name string `json:"name"`
		}{
			Name: account.Bank.Name,
		},
	}, nil
}

type PaystackDVAResponse struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Bank          struct {
		Name string `json:"name"`
	} `json:"bank"`
}

func (r *repository) GetDVABalance(ctx context.Context, accountNumber string) (float64, error) {
	// Get PayStack DVA balance
	paystackBalance, err := utils.GetPaystackDVABalance(accountNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to get DVA balance: %v", err)
	}

	// Get account details from Paystack
	paystackAccount, err := r.getPaystackDVAAccount(accountNumber)
	if err != nil {
		// If we can't get Paystack account details, just return Paystack balance
		return paystackBalance, nil
	}

	// First find the user by their full name
	var userID uint32
	if err := r.db.Table("users").
		Select("id").
		Where("fullname = ?", paystackAccount.AccountName).
		First(&userID).Error; err != nil {
		// If we can't find user, just return Paystack balance
		return paystackBalance, nil
	}

	// Then find the store associated with this user
	var store Store
	if err := r.db.Where("user_id = ?", userID).First(&store).Error; err != nil {
		// If we can't find store, just return Paystack balance
		return paystackBalance, nil
	}

	// Get store earnings for this specific store
	earnings, err := r.GetStoreEarnings(ctx, store.ID)
	if err != nil {
		// If error getting earnings, just return Paystack balance
		return paystackBalance, nil
	}

	// Calculate total earnings
	var totalEarnings float64
	for _, earning := range earnings {
		if earning.Status == "released" {
			totalEarnings += earning.Amount
		}
	}

	// Calculate total balance
	totalBalance := paystackBalance + totalEarnings

	// Update store's wallet column with the total balance
	if err := r.db.Model(&Store{}).
		Where("id = ?", store.ID).
		Update("wallet", totalBalance).Error; err != nil {
		return totalBalance, fmt.Errorf("failed to update store wallet: %v", err)
	}

	// Return combined balance
	return totalBalance, nil
}

func (r *repository) GetOrderByUUID(ctx context.Context, uuid string) (*Order, error) {
	var order Order

	// First get the basic order details
	if err := r.db.Where("uuid = ?", uuid).First(&order).Error; err != nil {
		return nil, err
	}

	// Unmarshal delivery details if present
	if order.DeliveryDetailsJSON != "" {
		var details DeliveryDetails
		if err := json.Unmarshal([]byte(order.DeliveryDetailsJSON), &details); err != nil {
			return nil, fmt.Errorf("failed to unmarshal delivery details: %v", err)
		}
		order.DeliveryDetails = &details
	}

	// Unmarshal customer details if present
	if order.CustomerJSON != "" {
		var customer Customer
		if err := json.Unmarshal([]byte(order.CustomerJSON), &customer); err != nil {
			return nil, fmt.Errorf("failed to unmarshal customer: %v", err)
		}
		order.Customer = &customer
	}

	// Get products for this order from the products table
	var products []TrackedProduct
	for _, storeID := range order.StoresID {
		var storeProducts []Product
		if err := r.db.Table("products").
			Where("store = ? AND deleted_at IS NULL", storeID).
			Find(&storeProducts).Error; err != nil {
			return nil, fmt.Errorf("failed to get products for store %s: %v", storeID, err)
		}

		// Convert each product to TrackedProduct
		for _, p := range storeProducts {
			products = append(products, TrackedProduct{
				ID:        p.ID,
				Name:      p.Name,
				Thumbnail: p.Thumbnail,
				Price:     p.Price,
				Store:     p.Store,
				Status:    "active",
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
			})
		}
	}
	order.Products = products

	return &order, nil
}

func (r *repository) UpdateProductUnitsSold(ctx context.Context, productID uint32) error {
	result := r.db.Model(&Product{}).Where("id = ?", productID).
		UpdateColumn("units_sold", gorm.Expr("units_sold + ?", 1))
	return result.Error
}

func (r *repository) GetAllStores(ctx context.Context, limit, offset int) ([]*Store, error) {
	var stores []*Store
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

// UpdateStoreBankDetails updates or creates bank details for a store
func (r *repository) UpdateStoreBankDetails(ctx context.Context, storeID uint32, account *WithdrawalAccount) error {
	// First, check how many accounts the store already has
	var count int64
	if err := r.db.Table("dva_accounts").Where("store_id = ?", storeID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count existing accounts: %v", err)
	}

	// Check if this account number already exists for this store
	var existingAccount struct {
		ID string
	}
	accountExists := r.db.Table("dva_accounts").
		Where("store_id = ? AND account_number = ?", storeID, account.AccountNumber).
		First(&existingAccount).Error == nil

	// If account exists, we'll update it. If not, check if we can add a new one
	if !accountExists && count >= 3 {
		return fmt.Errorf("maximum number of bank accounts (3) reached for this store")
	}

	// Generate a UUID for customer_id if not exists
	customerID := uuid.New()

	// Create or update the account record
	accountData := map[string]interface{}{
		"store_id":       storeID,
		"customer_id":    customerID,
		"bank_id":        account.BankID,
		"account_number": account.AccountNumber,
		"account_name":   account.AccountName,
		"bank_code":      account.BankCode,
		"bank_name":      account.BankName,
		"bank_image":     account.BankImage,
	}

	if accountExists {
		// Update existing account
		if err := r.db.Table("dva_accounts").
			Where("store_id = ? AND account_number = ?", storeID, account.AccountNumber).
			Updates(accountData).Error; err != nil {
			return fmt.Errorf("failed to update bank details: %v", err)
		}
	} else {
		// Create new account
		if err := r.db.Table("dva_accounts").Create(accountData).Error; err != nil {
			return fmt.Errorf("failed to create bank details: %v", err)
		}
	}

	return nil
}

func (r *repository) AddStoreEarnings(ctx context.Context, earnings *StoreEarnings) error {
	// Log the earnings object before creating
	log.Printf("AddStoreEarnings: Attempting to create earnings with StoreID=%d, OrderID=%s, Amount=%f, Status=%s", earnings.StoreID, earnings.OrderID, earnings.Amount, earnings.Status)

	if err := r.db.WithContext(ctx).Create(earnings).Error; err != nil {
		// Log the exact error from the Create operation
		log.Printf("AddStoreEarnings: Failed to create earnings: %v", err)
		return fmt.Errorf("failed to add store earnings: %v", err)
	}
	// Log success
	log.Printf("AddStoreEarnings: Successfully created earnings with ID=%d", earnings.ID)

	return nil
}

func (r *repository) GetStoreEarnings(ctx context.Context, storeID uint32) ([]*StoreEarnings, error) {
	var earnings []*StoreEarnings
	if err := r.db.WithContext(ctx).
		Where("store_id = ? AND status = ?", storeID, "released").
		Find(&earnings).Error; err != nil {
		return nil, fmt.Errorf("failed to get store earnings: %v", err)
	}
	return earnings, nil
}

func (r *repository) GetAllOrders(ctx context.Context) ([]*Order, error) {
	var orders []*Order
	if err := r.db.Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}
	return orders, nil
}

func (r *repository) CheckStoreEarningsDiscrepancy(ctx context.Context, storeID uint32) (int, float64, error) {
	// Get count of delivered orders for this store
	var deliveredOrdersCount int64
	if err := r.db.Model(&Order{}).
		Where("? = ANY(stores_id) AND status = ?", storeID, "delivered").
		Count(&deliveredOrdersCount).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count delivered orders: %v", err)
	}

	// Get total store earnings
	var totalEarnings float64
	if err := r.db.Model(&StoreEarnings{}).
		Where("store_id = ? AND status = ?", storeID, "released").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalEarnings).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to get total earnings: %v", err)
	}

	return int(deliveredOrdersCount), totalEarnings, nil
}

func (r *repository) UpdateWallet(ctx context.Context, storeID uint32, amount float64) error {
	result := r.db.Model(&Store{}).
		Where("id = ?", storeID).
		Update("wallet", gorm.Expr("wallet + ?", amount))

	if result.Error != nil {
		return fmt.Errorf("failed to update store wallet: %v", result.Error)
	}
	return nil
}

func (r *repository) UpdatePaystackBalance(ctx context.Context, storeID uint32, amount float64) error {
	result := r.db.Model(&Store{}).
		Where("id = ?", storeID).
		Update("paystack_balance", gorm.Expr("paystack_balance + ?", amount))

	if result.Error != nil {
		return fmt.Errorf("failed to update store Paystack balance: %v", result.Error)
	}
	return nil
}

// CreatePaystackDVAAccount creates and stores a Paystack DVA account for a store
func (r *repository) CreatePaystackDVAAccount(ctx context.Context, storeID uint32, account *PaystackDVAResponse, email string) error {
	// Generate a unique ID for the DVA account
	timestamp := time.Now().Unix()
	dvaID := fmt.Sprintf("PDVA_%d_%s", timestamp, utils.GenerateRandomString(8))

	// Create the DVA account record
	dvaAccount := map[string]interface{}{
		"id":             dvaID,
		"store_id":       storeID,
		"account_number": account.AccountNumber,
		"account_name":   account.AccountName,
		"bank_name":      account.Bank.Name,
		"email":          email,
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
	}

	// Insert the record
	if err := r.db.Table("paystack_dva_accounts").Create(dvaAccount).Error; err != nil {
		return fmt.Errorf("failed to create Paystack DVA account: %v", err)
	}

	return nil
}

// GetPaystackDVAAccount retrieves a store's Paystack DVA account
func (r *repository) GetPaystackDVAAccount(ctx context.Context, storeID uint32) (*PaystackDVAResponse, error) {
	var dvaAccount struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankName      string `json:"bank_name"`
		Email         string `json:"email"`
	}

	if err := r.db.Table("paystack_dva_accounts").
		Where("store_id = ?", storeID).
		First(&dvaAccount).Error; err != nil {
		return nil, fmt.Errorf("failed to get Paystack DVA account: %v", err)
	}

	return &PaystackDVAResponse{
		AccountNumber: dvaAccount.AccountNumber,
		AccountName:   dvaAccount.AccountName,
		Bank: struct {
			Name string `json:"name"`
		}{
			Name: dvaAccount.BankName,
		},
	}, nil
}

// DeleteExistingPaystackDVAAccounts deletes all existing Paystack DVA accounts from the database
func (r *repository) DeleteExistingPaystackDVAAccounts(ctx context.Context) error {
	// Use a WHERE clause with a condition that's always true
	if err := r.db.Table("paystack_dva_accounts").Where("1=1").Delete(&struct{}{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing Paystack DVA accounts: %v", err)
	}
	return nil
}

// SyncExistingPaystackDVAAccounts retrieves all existing Paystack DVA accounts and stores them in our database
func (r *repository) SyncExistingPaystackDVAAccounts(ctx context.Context) error {
	// First delete existing accounts
	if err := r.DeleteExistingPaystackDVAAccounts(ctx); err != nil {
		return fmt.Errorf("failed to delete existing accounts: %v", err)
	}

	// Get all users with type "seller" from our database
	var users []struct {
		ID    uint32 `gorm:"column:id"`
		Email string `gorm:"column:email"`
	}
	if err := r.db.Table("users").
		Where("usertype = ?", "seller").
		Select("id, email").
		Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch seller users: %v", err)
	}

	// For each seller user, get their store and create DVA account
	for _, user := range users {
		// Skip if user has no email
		if user.Email == "" {
			continue
		}

		// Get the store ID for this user
		var store struct {
			ID uint32 `gorm:"column:id"`
		}
		if err := r.db.Table("stores").
			Where("user_id = ?", user.ID).
			Select("id").
			First(&store).Error; err != nil {
			log.Printf("Warning: Failed to get store for user %d: %v", user.ID, err)
			continue
		}

		// Get account from Paystack API directly
		paystackAccount, err := r.getPaystackDVAAccount(user.Email)
		if err != nil {
			log.Printf("Warning: Failed to get Paystack DVA account for user %d: %v", user.ID, err)
			continue
		}

		// Generate a unique ID for the DVA account
		timestamp := time.Now().Unix()
		dvaID := fmt.Sprintf("PDVA_%d_%s", timestamp, utils.GenerateRandomString(8))

		// Create new account record
		paystack_dva_accounts := &PaystackDVAAccount{
			ID:            dvaID,
			StoreID:       store.ID,
			AccountNumber: paystackAccount.AccountNumber,
			AccountName:   paystackAccount.AccountName,
			BankName:      paystackAccount.Bank.Name,
			Email:         user.Email,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Save to database
		if err := r.db.Create(paystack_dva_accounts).Error; err != nil {
			log.Printf("Warning: Failed to save DVA account for user %d: %v", user.ID, err)
			continue
		}

		log.Printf("Successfully synced DVA account for user %d (store %d)", user.ID, store.ID)
	}

	return nil
}

type PaystackDVAAccount struct {
	ID            string    `gorm:"column:id;primaryKey"`
	StoreID       uint32    `gorm:"column:store_id"`
	AccountNumber string    `gorm:"column:account_number"`
	AccountName   string    `gorm:"column:account_name"`
	BankName      string    `gorm:"column:bank_name"`
	BankCode      string    `gorm:"column:bank_code"`
	Email         string    `gorm:"column:email"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

// TableName specifies the table name for PaystackDVAAccount
func (PaystackDVAAccount) TableName() string {
	return "paystack_dva_accounts"
}

func (r *repository) GetDB() *gorm.DB {
	return r.db
}

// UpdateWalletBalance updates the store's wallet with the total balance (Paystack + earnings)
func (r *repository) UpdateWalletBalance(ctx context.Context, storeID uint32) error {
	store, err := r.GetStore(ctx, storeID)
	if err != nil {
		return fmt.Errorf("failed to get store: %v", err)
	}

	// Get store earnings
	var totalEarnings float64
	var earnings []*StoreEarnings
	if err := r.db.Where("store_id = ? AND status = ?", storeID, "released").Find(&earnings).Error; err != nil {
		return fmt.Errorf("failed to get store earnings: %v", err)
	}

	// Calculate total earnings
	for _, earning := range earnings {
		totalEarnings += earning.Amount
	}

	// Get total withdrawals
	var totalWithdrawals float64
	var withdrawals []*Withdrawal
	if err := r.db.Where("store_id = ? AND status IN ?", storeID, []string{"pending", "approved", "completed"}).Find(&withdrawals).Error; err != nil {
		return fmt.Errorf("failed to get withdrawals: %v", err)
	}

	// Calculate total withdrawals
	for _, withdrawal := range withdrawals {
		totalWithdrawals += withdrawal.Amount
	}

	// Calculate total balance
	totalBalance := store.PaystackBalance + totalEarnings - totalWithdrawals

	// Update wallet if different
	if store.Wallet != totalBalance {
		if err := r.UpdateWallet(ctx, storeID, totalBalance-store.Wallet); err != nil {
			return fmt.Errorf("failed to update wallet: %v", err)
		}
	}

	return nil
}
