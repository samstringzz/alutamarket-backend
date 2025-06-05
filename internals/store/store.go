package store

import (
	"context"
	"time"

	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Transactions struct {
	gorm.Model
	StoreID   string    `json:"store_id" db:"store_id"`
	Status    string    `json:"status" db:"status"`
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

type InvoiceCustomer struct {
	Email  string `json:"email" db:"email"`
	Name   string `json:"name" db:"name"`
	Number string `json:"number" db:"number"`
}

type InvoiceItem struct {
	Quantity int32   `json:"quantity" db:"quantity"`
	Name     string  `json:"name" db:"name"`
	Price    float64 `json:"price" db:"price"`
}

type InvoiceDelivery struct {
	Option  string  `json:"option" db:"option"`
	Address string  `json:"address" db:"address"`
	Fee     float64 `json:"fee" db:"fee"`
}
type Invoice struct {
	gorm.Model
	Customer        *InvoiceCustomer `gorm:"serializer:json" json:"customer" db:"customer"`
	DueDate         string           `json:"due_date" db:"due_date"`
	Items           []*InvoiceItem   `gorm:"serializer:json"  json:"items" db:"items"`
	DeliveryDetails *InvoiceDelivery `gorm:"serializer:json" json:"delivery_details" db:"delivery_details"`
	StoreID         uint32           `json:"store_id" db:"store_id"`
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
	BankID        int    `json:"bank_id"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	BankImage     string `json:"bank_image"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}
type Store struct {
	gorm.Model
	ID                 uint32               `gorm:"primaryKey;uniqueIndex;not null;autoIncrement"  json:"id" db:"id"`
	Name               string               `json:"name" db:"name"`
	Customers          []Customer           `json:"customers" gorm:"foreignKey:ID"`
	UserID             uint32               `json:"user_id" db:"user_id"`
	Link               string               `json:"link" db:"link"`
	Description        string               `json:"description" db:"description"`
	HasPhysicalAddress bool                 `json:"hasphysical_address" db:"has_physical_address"`
	Address            string               `json:"address" db:"address"`
	Transactions       []*Transactions      `gorm:"serializer:json"`
	Followers          []*Follower          `gorm:"serializer:json"`
	Products           []Product            `gorm:"serializer:json"`
	Wallet             float64              `json:"wallet" db:"wallet"`
	PaystackBalance    float64              `json:"paystack_balance" db:"paystack_balance"`
	Status             bool                 `json:"status" db:"status"`
	Thumbnail          string               `json:"thumbnail" db:"thumbnail"`
	Phone              string               `json:"phone" db:"phone"`
	Email              string               `json:"email" db:"email"`
	Background         string               `json:"background" db:"background"`
	Visitors           pq.StringArray       `gorm:"type:text[]" json:"visitors"`
	Accounts           []*WithdrawalAccount `gorm:"serializer:json" json:"accounts" db:"accounts"`
	Orders             []*StoreOrder        `gorm:"serializer:json"`
	MaintenanceMode    bool                 `json:"maintenance_mode" db:"maintenance_mode"`
}

type UpdateStore struct {
	gorm.Model
	ID                 uint32             `gorm:"primaryKey;uniqueIndex;not null;autoIncrement"  json:"id" db:"id"`
	Name               string             `json:"name" db:"name"`
	UserID             uint32             `json:"user_id" db:"user_id"`
	Link               string             `json:"link" db:"link"`
	Description        string             `json:"description" db:"description"`
	HasPhysicalAddress bool               `json:"hasphysical_address" db:"has_physical_address"`
	Address            string             `json:"address" db:"address"`
	Transactions       []*Transactions    `gorm:"serializer:json"`
	Followers          []Follower         `gorm:"serializer:json"`
	Orders             []*StoreOrder      `gorm:"serializer:json"`
	Reviews            []*Review          `gorm:"serializer:json"`
	Account            *WithdrawalAccount `gorm:"serializer:json"`
	Products           []Product          `gorm:"serializer:json"`
	Wallet             float64            `json:"wallet" db:"wallet"`
	PaystackBalance    float64            `json:"paystack_balance" db:"paystack_balance"`
	Status             bool               `json:"status" db:"status"`
	Thumbnail          string             `json:"thumbnail" db:"thumbnail"`
	Phone              string             `json:"phone" db:"phone"`
	Email              string             `json:"email" db:"email"`
	Background         string             `json:"background" db:"background"`
	Visitors           []string           `gorm:"serializer:json" json:"visitors" db:"visitors"`
	MaintenanceMode    bool               `json:"maintenance_mode" db:"maintenance_mode"`
}

type TrackedProduct struct {
	gorm.Model
	ID        uint32    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Thumbnail string    `json:"thumbnail" db:"thumbnail"`
	Price     float64   `json:"price" db:"price"`
	File      *string   `json:"file" db:"file"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Store     string    `json:"store" db:"store"`
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
	CartID              uint32           `json:"cart_id" db:"cart_id"`
	Coupon              string           `json:"coupon,omitempty" db:"coupon"`
	Fee                 string           `json:"fee" db:"fee"`
	Status              string           `json:"status" db:"status"`
	UserID              string           `json:"user_id" db:"user_id" gorm:"index"`
	CustomerJSON        string           `json:"customer" db:"customer" gorm:"column:customer"`
	Customer            *Customer        `json:"-" gorm:"-"`
	SellerID            string           `json:"seller_id" db:"seller_id"`
	StoresID            pq.StringArray   `gorm:"type:text[]" json:"stores_id" db:"stores_id"`
	DeliveryDetailsJSON string           `json:"delivery_details" db:"delivery_details" gorm:"column:delivery_details"`
	DeliveryDetails     *DeliveryDetails `json:"-" gorm:"-"`
	Amount              string           `json:"amount" db:"amount"`
	UUID                string           `json:"uuid" db:"uuid"`
	PaymentGateway      string           `json:"payment_gateway" db:"payment_gateway"`
	PaymentMethod       string           `json:"payment_method" db:"payment_method"`
	TransRef            string           `json:"trans_ref" db:"trans_ref"`
	TransStatus         string           `json:"trans_status" db:"trans_status"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
	Products            []TrackedProduct `gorm:"serializer:json" json:"products" db:"products"`
}

type Customer struct {
	ID      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
	Info    string `json:"info" db:"info"`
	Email   string `json:"email" db:"email"`
}

type StoreProduct struct {
	Name      string  `json:"name" db:"name"`
	Thumbnail string  `json:"thumbnail" db:"thumbnail"`
	Status    string  `json:"status" db:"status"`
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
	Customer  Customer        `gorm:"serializer:json"  json:"customer" db:"customer"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}
type Buyer struct {
	Nickname string `json:"nickname" db:"nickname"`
	Avatar   string `json:"avatar" db:"avatar"`
	Comment  string `json:"comment" db:"comment"`
}

type Review struct {
	gorm.Model
	StoreID   uint32    `json:"store_id" db:"store_id"`
	ProductID uint32    `json:"product_id" db:"product_id"`
	OrderID   string    `json:"uuid" db:"uuid"`
	Buyer     *Buyer    `gorm:"serializer:json" json:"buyer" db:"buyer"`
	SellerID  uint32    `json:"seller_id" db:"seller_id"`
	Rating    float64   `json:"rating" db:"rating"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Fund struct {
	StoreID       uint32  `json:"store_id" db:"store_id"`
	UserID        uint32  `json:"user_id" db:"user_id"`
	Amount        float32 `json:"amount" db:"amount"`
	Email         string  `json:"email" db:"email"`
	AccountNumber string  `json:"account_number" db:"account_number"`
	BankCode      string  `json:"bank_code" db:"bank_code"`
}

type UpdateStoreOrderInput struct {
	StoreID uint32 `json:"store_id" db:"store_id"`
	UUID    string `json:"id" db:"id"`
	Status  string `json:"status" db:"status"`
}

type PaystackBank struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Add these structs after the PaymentDetails struct
type DVAAccount struct {
	ID            string      `json:"id" gorm:"type:varchar(100);primary_key"`
	AccountName   string      `json:"account_name" gorm:"type:varchar(100)"`
	AccountNumber string      `json:"account_number" gorm:"type:varchar(20)"`
	CustomerID    string      `json:"customer_id" gorm:"type:varchar(100)"`
	BankID        string      `json:"bank_id" gorm:"type:varchar(100)"`
	Customer      DVACustomer `json:"customer" gorm:"foreignKey:CustomerID"`
	Bank          DVABank     `json:"bank" gorm:"foreignKey:BankID"`
}

type DVACustomer struct {
	ID    string `json:"id" gorm:"type:varchar(100);primary_key"`
	Email string `json:"email" gorm:"type:varchar(100)"`
}

type DVABank struct {
	ID   string `json:"id" gorm:"type:varchar(100);primary_key"`
	Name string `json:"name" gorm:"type:varchar(100)"`
	Slug string `json:"slug" gorm:"type:varchar(50)"`
}

type StoreEarnings struct {
	gorm.Model
	StoreID         uint32    `json:"store_id" db:"store_id"`
	OrderID         string    `json:"order_id" db:"order_id"`
	Amount          float64   `json:"amount" db:"amount"`
	Status          string    `json:"status" db:"status"`                     // pending/released
	TransactionType string    `json:"transaction_type" db:"transaction_type"` // order/direct_transfer
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type WithdrawalStatus string

const (
	StatusPending   WithdrawalStatus = "pending"
	StatusApproved  WithdrawalStatus = "approved"
	StatusRejected  WithdrawalStatus = "rejected"
	StatusCompleted WithdrawalStatus = "completed"
)

type Withdrawal struct {
	ID                 uint32     `gorm:"primaryKey" json:"id"`
	StoreID            uint32     `json:"store_id"`
	Amount             float64    `json:"amount"`
	Status             string     `json:"status"`
	BankName           string     `json:"bank_name"`
	AccountNumber      string     `json:"account_number"`
	AccountName        string     `json:"account_name"`
	RejectionReason    string     `json:"rejection_reason,omitempty"`
	PaystackTransferID string     `json:"paystack_transfer_id,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type NewWithdrawal struct {
	StoreID       uint32  `json:"store_id"`
	Amount        float64 `json:"amount"`
	BankName      string  `json:"bank_name"`
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
}

type Repository interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	CreateInvoice(ctx context.Context, req *Invoice) (*Invoice, error)
	DeleteStore(ctx context.Context, id uint32) error
	CheckStoreName(ctx context.Context, query string) error
	UpdateStore(ctx context.Context, req *UpdateStore) (*Store, error)
	GetStore(ctx context.Context, id uint32) (*Store, error)
	GetStoreByName(ctx context.Context, name string) (*Store, error)
	UpdateOrderStatus(ctx context.Context, uuid string, status, transStatus string) error
	CreateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetOrder(ctx context.Context, storeId uint32, orderId string) (*Order, error)
	GetOrders(ctx context.Context, storeId uint32) ([]*Order, error)
	GetOrdersByStore(ctx context.Context, storeName string) ([]*Order, error)
	GetDVAAccount(ctx context.Context, email string) (*DVAAccount, error)
	GetDVABalance(ctx context.Context, id string) (float64, error)
	GetPurchasedOrders(ctx context.Context, userId string) ([]*Order, error)
	UpdateOrder(ctx context.Context, req *UpdateStoreOrderInput) (*Order, error)
	GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error)
	UpdateStoreFollowership(ctx context.Context, storeID uint32, follower *Follower, action string) (*Store, error)
	CreateTransactions(ctx context.Context, req *Transactions) (*Transactions, error)
	WithdrawFund(ctx context.Context, req *Fund) error
	GetInvoices(ctx context.Context, storeID uint32) ([]*Invoice, error)
	GetFollowedStores(ctx context.Context, userID uint32) ([]*Store, error)
	AddReview(ctx context.Context, review *Review) error
	GetReviews(ctx context.Context, filter string, value interface{}) ([]*Review, error)
	GetOrderByUUID(ctx context.Context, uuid string) (*Order, error)
	UpdateProductUnitsSold(ctx context.Context, productID uint32) error
	GetAllStores(ctx context.Context, limit, offset int) ([]*Store, error)
	UpdateStoreBankDetails(ctx context.Context, storeID uint32, account *WithdrawalAccount) error
	AddStoreEarnings(ctx context.Context, earnings *StoreEarnings) error
	GetStoreEarnings(ctx context.Context, storeID uint32) ([]*StoreEarnings, error)
	GetAllOrders(ctx context.Context) ([]*Order, error)
	CheckStoreEarningsDiscrepancy(ctx context.Context, storeID uint32) (int, float64, error)
	CreatePaystackDVAAccount(ctx context.Context, storeID uint32, account *PaystackDVAResponse, email string) error
	GetPaystackDVAAccount(ctx context.Context, storeID uint32) (*PaystackDVAResponse, error)
	SyncExistingPaystackDVAAccounts(ctx context.Context) error
	UpdatePaystackBalance(ctx context.Context, storeID uint32, amount float64) error
	UpdateWallet(ctx context.Context, storeID uint32, amount float64) error
	GetDB() *gorm.DB
}

type Service interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	CreateInvoice(ctx context.Context, req *Invoice) (*Invoice, error)
	UpdateStore(ctx context.Context, req *UpdateStore) (*Store, error)
	DeleteStore(ctx context.Context, id uint32) error
	GetStoreByName(ctx context.Context, name string) (*Store, error)
	CheckStoreName(ctx context.Context, query string) error
	GetStore(ctx context.Context, id uint32) (*Store, error)
	GetOrdersByStore(ctx context.Context, storeName string) ([]*Order, error)
	GetPurchasedOrders(ctx context.Context, userId string) ([]*Order, error)
	UpdateOrderStatus(ctx context.Context, uuid string, status, transStatus string) error
	CreateOrder(ctx context.Context, req *StoreOrder) (*StoreOrder, error)
	GetOrders(ctx context.Context, storeId uint32) ([]*Order, error)
	UpdateOrder(ctx context.Context, req *UpdateStoreOrderInput) (*Order, error)
	GetStores(ctx context.Context, user uint32, limit, offset int) ([]*Store, error)
	CreateTransactions(ctx context.Context, req *Transactions) (*Transactions, error)
	UpdateStoreFollowership(ctx context.Context, storeID uint32, follower *Follower, action string) (*Store, error)
	WithdrawFund(ctx context.Context, req *Fund) error
	GetInvoices(ctx context.Context, storeID uint32) ([]*Invoice, error)
	AddReview(ctx context.Context, review *Review) error
	GetReviews(ctx context.Context, filter string, value interface{}) ([]*Review, error)
	GetDVAAccount(ctx context.Context, email string) (*DVAAccount, error)
	GetDVABalance(ctx context.Context, id string) (float64, error)
	GetFollowedStores(ctx context.Context, userID uint32) ([]*Store, error)
	GetOrderByUUID(ctx context.Context, uuid string) (*Order, error)
	UpdateProductUnitsSold(ctx context.Context, productID uint32) error
	GetAllOrders(ctx context.Context) ([]*Order, error)
	AddStoreEarnings(ctx context.Context, earnings *StoreEarnings) error
	CheckStoreEarningsDiscrepancy(ctx context.Context, storeID uint32) (int, float64, error)
	CreatePaystackDVAAccount(ctx context.Context, storeID uint32, account *PaystackDVAResponse, email string) error
	GetPaystackDVAAccount(ctx context.Context, storeID uint32) (*PaystackDVAResponse, error)
	SyncExistingPaystackDVAAccounts(ctx context.Context) error
	GetStoreEarnings(ctx context.Context, storeID uint32) ([]*StoreEarnings, error)
}
