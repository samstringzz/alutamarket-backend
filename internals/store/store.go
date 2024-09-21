package store

import (
	"context"
	"time"

	"github.com/Chrisentech/aluta-market-api/internals/product"
	"gorm.io/gorm"
)

type Transactions struct {
	gorm.Model
	StoreID   string    `json:"store_id" db:"store_id"`
	Status    string    `json:"status" db:"status"` //pending,approved,canceled
	User      string    `json:"user" db:"user"`
	Amount    float64   `json:"amount" db:"amount"`
	UUID      string    `json:"uuid" db:"uuid"`
	Type      string    `json:"type" db:"type"`
	Category  string    `json:"category" db:"category"` // inovice/ transaction
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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

type Downloads struct {
	gorm.Model
	ID        string    `json:"id" db:"id"`
	Thumbnail string    `json:"thumbnail" db:"thumbnail"`
	Name      string    `json:"name" db:"name"`
	Price     float64   `json:"price" db:"price"`
	Discount  int       `json:"discount" db:"discount"`
	UUID      string    `json:"uuid" db:"uuid"`
	File      string    `json:"file" db:"file"`
	Users     []string  `gorm:"serializer:json" json:"paid_users" db:"paid_users"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WithdrawalAccount struct {
	BankName      string `json:"bank_name" db:"bank_name"`
	BankCode      string `json:"bank_code" db:"bank_code"`
	BankImage     string `json:"bank_image" db:"bank_image"`
	AccountNumber string `json:"account_number" db:"account_number"`
	AccountName   string `json:"account_name" db:"account_name"`
}
type Store struct {
	gorm.Model
	ID                 uint32               `gorm:"primaryKey;uniqueIndex;not null;autoIncrement"  json:"id" db:"id"`
	Name               string               `json:"name" db:"name"`
	UserID             uint32               `json:"user_id" db:"user_id"`
	Link               string               `json:"link" db:"link"`
	Description        string               `json:"description" db:"description"`
	HasPhysicalAddress bool                 `json:"hasphysical_address" db:"has_physical_address"`
	Address            string               `json:"address" db:"address"`
	Transactions       []*Transactions      `gorm:"serializer:json"`
	Followers          []*Follower          `gorm:"many2many:store_followers;" json:"store_followers"`
	Orders             []*StoreOrder        `gorm:"serializer:json"`
	Products           []Product            `gorm:"serializer:json"`
	Wallet             float64              `json:"wallet" db:"wallet"`
	Status             bool                 `json:"status" db:"status"`
	Thumbnail          string               `json:"thumbnail" db:"thumbnail"`
	Phone              string               `json:"phone" db:"phone"`
	Email              string               `json:"email" db:"email"`
	Background         string               `json:"background" db:"background"`
	Visitors           []string             `gorm:"serializer:json" json:"visitors" db:"visitors"`
	Accounts           []*WithdrawalAccount `gorm:"serializer:json" json:"accounts" db:"accounts"`
}

type UpdateStore struct {
	gorm.Model
	ID                 uint32          `gorm:"primaryKey;uniqueIndex;not null;autoIncrement"  json:"id" db:"id"`
	Name               string          `json:"name" db:"name"`
	UserID             uint32          `json:"user_id" db:"user_id"`
	Link               string          `json:"link" db:"link"`
	Description        string          `json:"description" db:"description"`
	HasPhysicalAddress bool            `json:"hasphysical_address" db:"has_physical_address"`
	Address            string          `json:"address" db:"address"`
	Transactions       []*Transactions `gorm:"serializer:json"`
	Followers          []Follower      `gorm:"serializer:json"`
	Orders             []*StoreOrder   `gorm:"serializer:json"`
	Products           []Product       `gorm:"serializer:json"`
	Wallet             float64         `json:"wallet" db:"wallet"`
	Status             bool            `json:"status" db:"status"`
	Thumbnail          string          `json:"thumbnail" db:"thumbnail"`
	Phone              string          `json:"phone" db:"phone"`
	Email              string          `json:"email" db:"email"`
	Background         string          `json:"background" db:"background"`
	Visitors           string          `gorm:"serializer:json" json:"visitors" db:"visitors"`
}

type TrackedProduct struct {
	gorm.Model
	ID        uint32    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Thumbnail string    `json:"thumbnail" db:"thumbnail"`
	Price     float64   `json:"price" db:"price"`
	File      *string   `json:"file" db:"file"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Discount  float64   `json:"discount" db:"discount"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
type DeliveryDetails struct {
	Method  string  `json:"method,omitempty" db:"method"`
	Address string  `json:"address,omitempty" db:"address"`
	Fee     float64 `json:"fee,omitempty" db:"fee"`
}

// Purchased Orders
type Order struct {
	gorm.Model
	// StoresID       string  `gorm:"serializer:json" json:"store" db:"store_id"`
	CartID          uint32           `json:"cart_id" db:"cart_id"`
	Coupon          string           `json:"coupon,omitempty" db:"coupon"`
	Fee             float64          `json:"fee" db:"fee"`
	Status          string           `json:"status" db:"status"` //order status
	UserID          string           `json:"user_id" db:"user_id"`
	DeliveryDetails DeliveryDetails  `gorm:"serializer:json" json:"delivery_details" db:"delivery_details"`
	Amount          float64          `json:"amount" db:"amount"`
	UUID            string           `json:"uuid" db:"uuid"`
	PaymentGateway  string           `json:"payment_gateway" db:"payment_gateway"`
	PaymentMethod   string           `json:"payment_method" db:"payment_method"`
	TransRef        string           `json:"trt_ref" db:"trt_ref"`
	TransStatus     string           `json:"txt_status" db:"txt_status"` //Transaction status
	Products        []TrackedProduct `gorm:"serializer:json" json:"products" db:"products"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" db:"updated_at"`
}

type Customer struct {
	ID      uint32 `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
	Info    string `json:"info" db:"info"`
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
	Products  []*StoreProduct `gorm:"serializer:json" json:"products" db:"products"`
	Status    string          `json:"status" db:"status"`
	TransRef  string          `json:"trt_ref" db:"trt_ref"`
	Active    bool            `json:"active" db:"active"`
	UUID      string          `json:"uuid" db:"uuid"`
	Customer  Customer        `json:"customer" db:"customer"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type Review struct {
}

type Fund struct {
	StoreID       uint32  `json:"store_id" db:"store_id"`
	UserID        uint32  `json:"user_id" db:"user_id"`
	Amount        float32 `json:"amount" db:"amount"`
	Email         string  `json:"email" db:"email"`
	AccountNumber string  `json:"account_number" db:"account_number"`
	BankCode      string  `json:"bank_code" db:"bank_code"`
}

type Repository interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	DeleteStore(ctx context.Context, id uint32) error
	CheckStoreName(ctx context.Context, query string) error
	UpdateStore(ctx context.Context, req *UpdateStore) (*Store, error)
	GetStore(ctx context.Context, id uint32) (*Store, error)
	GetStoreByName(ctx context.Context, name string) (*Store, error)
	CreateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetOrder(ctx context.Context, storeId uint32, orderId string) (*StoreOrder, error)
	GetOrders(ctx context.Context, storeId uint32) ([]*StoreOrder, error)
	GetPurchasedOrders(ctx context.Context, userId string) ([]*Order, error)
	UpdateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error)
	UpdateStoreFollowership(ctx context.Context, storeID uint32, follower *Follower, action string) (*Store, error)
	CreateTransactions(ctx context.Context, req *Transactions) (*Transactions, error)
	WithdrawFund(ctx context.Context, req *Fund) error
}

type Service interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	UpdateStore(ctx context.Context, req *UpdateStore) (*Store, error)
	DeleteStore(ctx context.Context, id uint32) error
	GetStoreByName(ctx context.Context, name string) (*Store, error)
	CheckStoreName(ctx context.Context, query string) error
	GetStore(ctx context.Context, id uint32) (*Store, error)
	GetPurchasedOrders(ctx context.Context, userId string) ([]*Order, error)
	CreateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetOrders(ctx context.Context, storeId uint32) ([]*StoreOrder, error)
	UpdateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error)
	CreateTransactions(ctx context.Context, req *Transactions) (*Transactions, error)
	UpdateStoreFollowership(ctx context.Context, storeID uint32, follower *Follower, action string) (*Store, error)
	WithdrawFund(ctx context.Context, req *Fund) error
}
