package store

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	dbURI := os.Getenv("DB_URI")

	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return &repository{
		db: db,
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
	if req.Visitors != "" {
		existingStore.Visitors = append(existingStore.Visitors, req.Visitors)
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

	// Convert []*StoreProduct to []StoreProduct
	var products []*StoreProduct
	products = append(products, req.Products...)

	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	req.UUID = utils.GenerateUUID()
	req.Products = products
	req.Status = "pending"
	store.Orders = append(store.Orders, req)

	err = r.db.WithContext(ctx).Save(&store).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *repository) GetOrders(ctx context.Context, storeID uint32) ([]*StoreOrder, error) {
	var store Store
	err := r.db.WithContext(ctx).Preload("Orders").First(&store, storeID).Error
	if err != nil {
		return nil, err
	}
	return store.Orders, nil
}

func (r *repository) GetPurchasedOrders(ctx context.Context, userID string) ([]*Order, error) {
	var orders []*Order
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *repository) GetOrder(ctx context.Context, storeID uint32, orderID string) (*StoreOrder, error) {
	var store Store
	err := r.db.WithContext(ctx).Preload("Orders", "id = ?", orderID).First(&store, storeID).Error
	if err != nil {
		return nil, err
	}

	if len(store.Orders) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return store.Orders[0], nil
}

func (r *repository) UpdateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error) {
	storeID, err := strconv.ParseUint(req.StoreID, 10, 16)
	if err != nil {
		return nil, err
	}
	existingStore, err := r.GetStore(ctx, uint32(storeID))
	if err != nil {
		return nil, err
	}
	if len(existingStore.Orders) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Find the order by UUID
	var orderToUpdate *StoreOrder
	for i, order := range existingStore.Orders {
		if order.UUID == req.UUID {
			orderToUpdate = existingStore.Orders[i]
			break
		}
	}

	if orderToUpdate == nil {
		return nil, gorm.ErrRecordNotFound
	}

	// Update the order fields
	orderToUpdate.Status = req.Status
	orderToUpdate.UpdatedAt = time.Now()

	// Save the changes to the store
	err = r.db.WithContext(ctx).Save(&existingStore).Error
	if err != nil {
		return nil, err
	}

	return orderToUpdate, nil
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
		return errors.NewAppError(http.StatusNotFound, "NOT FOUND", "Oops, An error occured in transaction")
	}
	if req.Amount > float32(store.Wallet) {
		return errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Your Wallet amount is not within range of withdrawal amount")
	}
	err = utils.PayFund(req.Amount, req.Email, req.AccountNumber, req.BankCode)
	if err != nil {
		return err
	}
	return nil
}
