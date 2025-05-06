package store

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Chrisentech/aluta-market-api/database"
	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/utils"
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

	newStore := &Store{
		Name:               req.Name,
		Link:               req.Link,
		HasPhysicalAddress: req.HasPhysicalAddress,
		UserID:             req.UserID,
		Wallet:             0,
		Address:            req.Address,
		Description:        req.Description,
		Status:             true,
		Phone:              req.Phone,
	}

	if err := r.db.Create(newStore).Error; err != nil {
		return nil, err
	}
	return newStore, nil
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

func (r *repository) GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error) {
	var stores []*Store
	if err := r.db.Where("user_id=?", user).Limit(limit).Offset(offset).Find(&stores).Error; err != nil {
		return nil, err
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
		existingStore.Visitors = append(existingStore.Visitors, req.Visitors...)
	}

	if req.Account != nil {
		existingStore.Accounts = append(existingStore.Accounts, req.Account)
	}
	if req.HasPhysicalAddress != existingStore.HasPhysicalAddress {
		existingStore.HasPhysicalAddress = req.HasPhysicalAddress
	}
	if req.Status != existingStore.Status {
		existingStore.Status = req.Status
	}
	if req.Address != "" {
		existingStore.Address = req.Address
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

	// Update the Store in the repository
	err = r.db.Save(existingStore).Error
	if err != nil {
		return nil, err
	}

	return existingStore, nil
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
	result := r.db.WithContext(ctx).Model(&Order{}).
		Where("uuid = ?", uuid).
		Updates(map[string]interface{}{
			"status":       status,
			"trans_status": transStatus,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no order found with UUID: %s", uuid)
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
	var store *Store
	err := r.db.First(&store, req.StoreID).Error
	if err != nil {
		return err
	}
	if req.UserID != store.UserID {
		return errors.NewAppError(http.StatusNotFound, "NOT_FOUND", "Oops, An error occurred in transaction")
	}
	if req.Amount > float32(store.Wallet) {
		return errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Your Wallet amount is not within range of withdrawal amount")
	}
	err = utils.PayFund(req.Amount, req.AccountNumber, req.BankCode)
	if err != nil {
		return err
	}
	return nil
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
	url := "https://api.paystack.co/dedicated_account"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("email", email)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("Paystack API error: %s", string(bodyBytes))
	}

	var response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			AccountNumber string `json:"account_number"`
			AccountName   string `json:"account_name"`
			Bank          struct {
				Name string `json:"name"`
			} `json:"bank"`
		} `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if !response.Status {
		return nil, fmt.Errorf("Paystack error: %s", response.Message)
	}

	// Check if there are any accounts returned
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no DVA account found for email: %s", email)
	}

	// Return the first account (assuming it's the most relevant one)
	account := &PaystackDVAResponse{
		AccountNumber: response.Data[0].AccountNumber,
		AccountName:   response.Data[0].AccountName,
		Bank: struct {
			Name string `json:"name"`
		}{
			Name: response.Data[0].Bank.Name,
		},
	}

	return account, nil
}

type PaystackDVAResponse struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Bank          struct {
		Name string `json:"name"`
	} `json:"bank"`
}

func (r *repository) GetDVABalance(ctx context.Context, accountNumber string) (float64, error) {
	url := "https://api.paystack.co/transaction"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, err
	}

	// Use virtual_account_number instead of recipient_account
	q := req.URL.Query()
	q.Add("virtual_account_number", accountNumber)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			Amount          float64 `json:"amount"`
			Status          string  `json:"status"`
			Currency        string  `json:"currency"`
			Channel         string  `json:"channel"`
			GatewayResponse string  `json:"gateway_response"`
			Metadata        struct {
				ReceiverAccountNumber string `json:"receiver_account_number"`
			} `json:"metadata"`
		} `json:"data"`
		Meta struct {
			TotalVolume float64 `json:"total_volume"`
		} `json:"meta"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}

	if !response.Status {
		return 0, fmt.Errorf("failed to get transactions: %s", response.Message)
	}

	// Use the total_volume from meta, which represents the total amount
	// Convert from kobo to naira
	balance := response.Meta.TotalVolume / 100

	return balance, nil
}
