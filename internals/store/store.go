package store

import (
	"context"
	"time"

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

type Product product.Product

type Follower struct {
	gorm.Model
	FollowerID    uint32 `json:"follower_id" db:"follower_id"`
	FollowerName  string `json:"follower" db:"follower_name"`
	StoreID       uint32 `json:"store" db:"store_id"`
	FollowerImage string `json:"follower_image" db:"follower_image"`
}

type DVADetails struct {
	UserID    string `json:"user_id" db:"user_id"`
	StoreName string `json:"store_name" db:"store_name"`
}

type Store struct {
	gorm.Model
	ID                 uint32        `gorm:"primaryKey;uniqueIndex;not null;autoIncrement"  json:"id" db:"id"`
	Name               string        `json:"name" db:"name"`
	UserID             uint32        `json:"user_id" db:"user_id"`
	Link               string        `json:"link" db:"link"`
	Description        string        `json:"description" db:"description"`
	HasPhysicalAddress bool          `json:"hasphysical_address" db:"has_physical_address"`
	Address            string        `json:"address" db:"address"`
	Transactions       Transactions  `gorm:"serializer:json"`
	Followers          []Follower    `gorm:"serializer:json"`
	Orders             []*StoreOrder `gorm:"serializer:json"`
	Products           []Product     `gorm:"serializer:json"`
	Wallet             float64       `json:"wallet" db:"wallet"`
	Status             bool          `json:"status" db:"status"`
	Thumbnail          string        `json:"thumbnail" db:"thumbnail"`
	Phone              string        `json:"phone" db:"phone"`
	Email              string        `json:"email" db:"email"`
	Background         string        `json:"background" db:"background"`
}

type TrackedProduct struct {
	gorm.Model
	ID        uint32    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Thumbnail string    `json:"thumbnail" db:"thumbnail"`
	Price     float64   `json:"price" db:"price"`
	Discount  float64   `json:"discount" db:"discount"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
type Order struct {
	gorm.Model
	// StoresID       string  `gorm:"serializer:json" json:"store" db:"store_id"`
	CartID         uint32           `json:"cart_id" db:"cart_id"`
	Coupon         string           `json:"coupon,omitempty" db:"coupon"`
	Fee            float64          `json:"fee" db:"fee"`
	Status         string           `json:"status" db:"status"` //pending,completed,failed
	UserID         string           `json:"user_id" db:"user_id"`
	Amount         float64          `json:"amount" db:"amount"`
	UUID           string           `json:"uuid" db:"uuid"`
	PaymentGateway string           `json:"payment_gateway" db:"payment_gateway"`
	Products       []TrackedProduct `gorm:"serializer:json" json:"products" db:"products"`
}

type Customer struct {
	ID      uint32 `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
}
type StoreProduct struct {
	Name      string  `json:"name" db:"name"`
	Thumbnail string  `json:"thumbnail" db:"thumbnail"`
	Price     float64 `json:"price" db:"price"`
	Quantity  int     `json:"quantity" db:"quantity"`
	ID        uint32  `json:"id" db:"id"`
}
type StoreOrder struct {
	StoreID   string          `gorm:"serializer:json" json:"store" db:"store_id"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
	Products  []*StoreProduct `gorm:"serializer:json" json:"products" db:"products"`
	Status    string          `json:"status" db:"status"`
	UUID      string          `json:"uuid" db:"uuid"`
	Customer  `json:"customer" db:"customer"`
}

type Review struct {
}

type Repository interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	DeleteStore(ctx context.Context, id uint32) error
	CheckStoreName(ctx context.Context, query string) error
	UpdateStore(ctx context.Context, req *Store) (*Store, error)
	GetStore(ctx context.Context, id uint32) (*Store, error)
	GetStoreByName(ctx context.Context, name string) (*Store, error)
	CreateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetOrder(ctx context.Context, storeId uint32, orderId string) (*StoreOrder, error)
	GetOrders(ctx context.Context, storeId uint32) ([]*StoreOrder, error)
	UpdateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error)
	UpdateStoreFollowership(ctx context.Context, storeID uint32, follower Follower, action string) (*Store, error)
}

type Service interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	UpdateStore(ctx context.Context, req *Store) (*Store, error)
	DeleteStore(ctx context.Context, id uint32) error
	GetStoreByName(ctx context.Context, name string) (*Store, error)
	CheckStoreName(ctx context.Context, query string) error
	GetStore(ctx context.Context, id uint32) (*Store, error)
	CreateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetOrders(ctx context.Context, storeId uint32) ([]*StoreOrder, error)
	UpdateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error)
	UpdateStoreFollowership(ctx context.Context, storeID uint32, follower Follower, action string) (*Store, error)
}
