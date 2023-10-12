package cart

import (
	"context"

	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"gorm.io/gorm"
)

type Product *product.Product
type User *user.User
type Order *store.Order
type Cart struct {
	gorm.Model
	ID     uint32       `gorm:"primaryKey;uniqueIndex;not null"`
	Items  []*CartItems `gorm:"serializer:json" json:"items" db:"items"`
	Total  float64      `json:"total" db:"total"`
	Active bool         `json:"active" db:"active"`
	UserID uint32       `json:"user" db:"user_id"`
}

type CartItems struct {
	gorm.Model
	Product  Product `gorm:"embedded"`
	CartID   uint32  `json:"cart" db:"cart_id"`
	Quantity int     `json:"quantity" db:"quantity"`
}

type Repository interface {
	ModifyCart(ctx context.Context, req *CartItems, user uint32) (*Cart, error)
	RemoveAllCart(ctx context.Context, id uint32) (*Cart, error)
	GetCart(ctx context.Context, user uint32) (*Cart, error)
	MakePayment(ctx context.Context, req Order) (*Order, error)
	InitiatePayment(ctx context.Context, req Order) (string, error)
}

type Service interface {
	ModifyCart(ctx context.Context, req *CartItems, user uint32) (*Cart, error)
	RemoveAllCart(ctx context.Context, id uint32) (*Cart, error)
	GetCart(ctx context.Context, user uint32) (*Cart, error)
	MakePayment(ctx context.Context, req Order) (*Order, error)
	InitiatePayment(ctx context.Context, req Order) (string, error)
}
