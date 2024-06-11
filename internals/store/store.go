package store

import (
	"context"

	"github.com/Chrisentech/aluta-market-api/internals/product"
	"gorm.io/gorm"
)

type Transactions struct {
	gorm.Model
	StoresID       []string `gorm:"serializer:json" json:"store" db:"store_id"`
	CartID         uint32   `json:"cart_id" db:"cart_id"`
	Coupon         string   `json:"coupon,omitempty" db:"coupon"`
	Fee            float64  `json:"fee" db:"fee"`
	Status         string   `json:"status" db:"status"` //pending,completed,failed
	UserID         string   `json:"user_id" db:"user_id"`
	Amount         float64  `json:"amount" db:"amount"`
	UUID           string   `json:"uuid" db:"uuid"`
	PaymentGateway string   `json:"payment_gateway" db:"payment_gateway"`
}

type Product *product.Product

type Follower struct {
	gorm.Model
	FollowerID    uint32 `json:"follower_id" db:"follower_id"`
	FollowerName  string `json:"follower" db:"follower_name"`
	StoreID       uint32 `json:"store" db:"store_id"`
	FollowerImage string `json:"follower_image" db:"follower_image"`
}

type Store struct {
	gorm.Model
	ID                 uint32       `gorm:"primaryKey;uniqueIndex;not null;autoIncrement"  json:"id" db:"id"`
	Name               string       `json:"name" db:"name"`
	UserID             uint32       `json:"user_id" db:"user_id"`
	Link               string       `json:"link" db:"link"`
	Description        string       `json:"description" db:"description"`
	HasPhysicalAddress bool         `json:"hasphysical_address" db:"has_physical_address"`
	Address            string       `json:"address" db:"address"`
	Transactions       Transactions `gorm:"serializer:json"`
	Followers          []Follower   `gorm:"serializer:json"`
	Orders             []Order      `gorm:"serializer:json"`
	Products           []Product    `gorm:"serializer:json"`
	Wallet             float64      `json:"wallet" db:"wallet"`
	Status             bool         `json:"status" db:"status"`
	Thumbnail          string       `json:"thumbnail" db:"thumbnail"`
	Phone              string       `json:"phone" db:"phone"`
	Email              string       `json:"email" db:"email"`
	Background         string       `json:"background" db:"background"`
}

type Order struct {
	gorm.Model
	StoresID       string  `gorm:"serializer:json" json:"store" db:"store_id"`
	CartID         uint32  `json:"cart_id" db:"cart_id"`
	Coupon         string  `json:"coupon,omitempty" db:"coupon"`
	Fee            float64 `json:"fee" db:"fee"`
	Status         string  `json:"status" db:"status"` //pending,completed,failed
	UserID         string  `json:"user_id" db:"user_id"`
	Amount         float64 `json:"amount" db:"amount"`
	UUID           string  `json:"uuid" db:"uuid"`
	PaymentGateway string  `json:"payment_gateway" db:"payment_gateway"`
}

type Repository interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	DeleteStore(ctx context.Context, id uint32) error
	UpdateStore(ctx context.Context, req *Store) (*Store, error)
	GetStore(ctx context.Context, id uint32) (*Store, error)
	GetStoreByName(ctx context.Context, name string) (*Store, error)
	GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error)
}

type Service interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	UpdateStore(ctx context.Context, req *Store) (*Store, error)
	DeleteStore(ctx context.Context, id uint32) error
	GetStoreByName(ctx context.Context, name string) (*Store, error)
	GetStore(ctx context.Context, id uint32) (*Store, error)
	GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error)
}
