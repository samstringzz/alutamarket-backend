package store

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Chrisentech/aluta-market-api/errors"
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

func (r *repository) CheckStoreName(ctx context.Context, query string)error{
	var stores []*Store
	if err := r.db.Where("name ILIKE ?","%"+query+"%").Find(&stores).Error; err != nil {
		return err
	}
	for _,item := range stores{
		if item.Name == query{
			return errors.NewAppError(http.StatusConflict, "CONFLICT", "Store Name already choosen")
		}
	}
	return nil
}

func (r *repository) CreateStore(ctx context.Context, req *Store) (*Store, error) {
	var count int64
	r.db.Model(&Store{}).Where("name =?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Store already exist")
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
	err := r.db.Where("name = ?", name).First(&store).Error
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

func (r *repository) UpdateStore(ctx context.Context, req *Store) (*Store, error) {

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
	if req.Link != "" {
		existingStore.Link = req.Link
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
	if req.Email != "" {
		existingStore.Email = req.Email
	}
	if req.Phone != "" {
		existingStore.Phone = req.Phone
	}
	if req.Background != "" {
		existingStore.Background = req.Background
	}
	existingStore.Wallet += req.Wallet

	// Update the Store in the repository
	err = r.db.Save(existingStore).Error
	if err != nil {
		return nil, err
	}

	return existingStore, nil
}

