package cart

import (
	"context"

	"github.com/Chrisentech/aluta-market-api/internals/product"
	"github.com/Chrisentech/aluta-market-api/internals/user"
	"gorm.io/gorm"
)

type Product *product.Product
type User *user.User
type Cart struct {
	gorm.Model
	ID     uint32         `gorm:"primaryKey;uniqueIndex;not null"`
	Items  []CartItems 
	Total  float64      `json:"total" db:"total"`
	Active bool         `json:"active" db:"active"`
	UserID uint32         `json:"user" db:"user_id"`
	User   User         `gorm:"foreignKey:UserID"`
}

type CartItems struct {
	gorm.Model
	Product  Product     `gorm:"embedded"`
	CartID uint32         `json:"cart" db:"cart_id"`
	Quantity uint32 `json:"quantity" db:"quantity"`
}

type Repository interface {
	AddToCart(ctx context.Context, req []*CartItems, user uint32) (*Cart, error)
	RemoveFromCart(ctx context.Context, id uint32) (*Cart, error)
	GetCart(ctx context.Context, user uint32) (*Cart, error)
}

type Service interface {
	AddToCart(ctx context.Context, req []*CartItems, user uint32) (*Cart, error)
	RemoveFromCart(ctx context.Context, id uint32) (*Cart, error)
	GetCart(ctx context.Context, user uint32) (*Cart, error)
}
